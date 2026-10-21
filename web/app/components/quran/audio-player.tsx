import { keepScreenAwake } from '~/utils/screen-wake-lock';
import { ContinuousAudio } from './continuous-audio';
import { AudioDiagnostics } from './audio-diagnostics';
import { Pause, Play, SkipBack, SkipForward, Volume2, VolumeX } from 'lucide-react';
import React, { useEffect, useRef, useState } from 'react';

export type AudioTrack = { url: string; language: 'arabic' | 'english'; label: string };
type AudioPlayerProps = {
  arabicUrl?: string;
  englishUrl?: string;
  verseLabel: string;
  playlist?: AudioTrack[];
};

export function AudioPlayer({ arabicUrl, englishUrl, verseLabel, playlist }: AudioPlayerProps) {
  const tracks = playlist || [
    ...(arabicUrl ? [{ url: arabicUrl, language: 'arabic' as const, label: verseLabel }] : []),
    ...(englishUrl ? [{ url: englishUrl, language: 'english' as const, label: verseLabel }] : []),
  ];
  const tracksRef = useRef(tracks);
  tracksRef.current = tracks;
  const queueKey = JSON.stringify(tracks);
  const audioRef = useRef<HTMLAudioElement>(null);
  const indexRef = useRef(-1);
  const streamRef = useRef<ContinuousAudio | null>(null);
  const [index, setIndex] = useState(-1);
  const [isPlaying, setIsPlaying] = useState(false);
  const [progress, setProgress] = useState(0);
  const [duration, setDuration] = useState(0);
  const [volume, setVolume] = useState(1);
  const [isMuted, setIsMuted] = useState(false);
  const [error, setError] = useState('');
  const activeTrack = tracks[index];
  const currentTrack = activeTrack?.language || null;

  useEffect(() => {
    if (isPlaying) return keepScreenAwake();
  }, [isPlaying]);

  const playbackChanged = (playing: boolean) => {
    setIsPlaying(playing);
    if ('mediaSession' in navigator) navigator.mediaSession.playbackState = playing ? 'playing' : 'paused';
  };

  const play = () => {
    const audio = audioRef.current;
    if (!audio) return;
    // Supported browsers can distinguish long-form playback from ambient sound.
    const session = (navigator as Navigator & { audioSession?: { type: string } }).audioSession;
    try { if (session) session.type = 'playback'; } catch { /* optional API */ }
    setError('');
    void audio.play().catch((reason: DOMException) => {
      if (reason.name === 'AbortError') return; // A newer track or pause superseded this request.
      playbackChanged(false);
      setError('Playback could not start. Press play to retry.');
    });
  };

  const updateTrack = (next: number) => {
    const audio = audioRef.current;
    const track = tracksRef.current[next];
    if (!audio || !track) return;
    indexRef.current = next;
    setIndex(next);
    setProgress(0);
    setDuration(0);

    if ('mediaSession' in navigator && typeof MediaMetadata !== 'undefined') {
      navigator.mediaSession.metadata = new MediaMetadata({
        title: track.label,
        artist: track.language === 'arabic' ? 'Arabic Recitation' : 'English Translation',
        album: 'Reminder',
        artwork: [{ src: '/icon-192.png', sizes: '192x192', type: 'image/png' }],
      });
    }
  };

  const startTrack = (next: number) => {
    const audio = audioRef.current;
    if (!audio || !tracksRef.current[next]) return;
    const continuous = new URLSearchParams(location.search).get('continuousAudio') === '1';
    if (continuous && !ContinuousAudio.supported()) {
      setError('Continuous audio is not supported by this browser. Remove continuousAudio=1 from the URL to use the standard player.');
      return;
    }
    updateTrack(next);
    if (continuous) {
      if (!streamRef.current?.seek(next)) {
        streamRef.current?.dispose();
        streamRef.current = new ContinuousAudio(audio, tracksRef.current.map(track => track.url), next,
          (segment, elapsed) => {
            if (indexRef.current !== segment.index) updateTrack(segment.index);
            setProgress(elapsed);
            const length = segment.end - segment.start;
            setDuration(length);
            try {
              navigator.mediaSession?.setPositionState?.({ duration: length, playbackRate: audio.playbackRate,
                position: Math.min(length, Math.max(0, elapsed)) });
            } catch { /* optional lock-screen seek support */ }
          }, message => { playbackChanged(false); setError(message); });
      }
    } else audio.src = tracksRef.current[next].url;
    // Start directly from the media event. No background timer or React render
    // is needed to advance to English or to the next verse.
    play();
  };

  const resume = () => {
    const audio = audioRef.current;
    if (!audio) return;
    if (streamRef.current?.hasFailed) startTrack(Math.max(0, indexRef.current));
    else if (indexRef.current < 0 || audio.ended) startTrack(0);
    else play();
  };
  const pause = () => audioRef.current?.pause();
  const togglePlay = () => audioRef.current?.paused ? resume() : pause();
  const skipLanguage = (language: AudioTrack['language']) => {
    const label = tracksRef.current[Math.max(indexRef.current, 0)]?.label;
    const next = tracksRef.current.findIndex(track => track.label === label && track.language === language);
    if (next >= 0) startTrack(next);
  };
  const skipToArabic = () => skipLanguage('arabic');
  const skipToEnglish = () => skipLanguage('english');
  const handleEnded = () => {
    if (streamRef.current) { playbackChanged(false); return; }
    if (indexRef.current + 1 < tracksRef.current.length) startTrack(indexRef.current + 1);
    else playbackChanged(false);
  };

  const seek = (value: number) => {
    if (streamRef.current) streamRef.current.seek(indexRef.current, value);
    else if (audioRef.current) audioRef.current.currentTime = value;
  };

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;
    audio.pause();
    indexRef.current = -1;
    setIndex(-1);
    setProgress(0);
    setDuration(0);
    setIsPlaying(false);
    setError('');
    // Prime the first resource without starting playback before a user gesture.
    if (tracksRef.current[0]) audio.src = tracksRef.current[0].url;
    else audio.removeAttribute('src');
    const session = 'mediaSession' in navigator ? navigator.mediaSession : null;
    const handler = (action: MediaSessionAction, callback: MediaSessionActionHandler | null) => {
      try { session?.setActionHandler(action, callback); } catch { /* unsupported action */ }
    };
    if (session) { session.playbackState = 'none'; session.metadata = null; }
    handler('play', resume);
    handler('pause', pause);
    handler('nexttrack', () => startTrack(indexRef.current + 1));
    handler('previoustrack', () => startTrack(Math.max(0, indexRef.current - 1)));
    handler('seekto', details => { if (details.seekTime !== undefined) seek(details.seekTime); });
    return () => {
      streamRef.current?.dispose();
      streamRef.current = null;
      audio.pause();
      audio.removeAttribute('src');
      audio.load();
      for (const action of ['play', 'pause', 'nexttrack', 'previoustrack', 'seekto'] as MediaSessionAction[]) handler(action, null);
      if (session) {
        session.playbackState = 'none'; session.metadata = null;
        try { session.setPositionState?.(); } catch { /* optional API */ }
      }
    };
  }, [queueKey]);

  useEffect(() => {
    if (audioRef.current) { audioRef.current.volume = volume; audioRef.current.muted = isMuted; }
  }, [volume, isMuted]);

  const toggleMute = () => setIsMuted(!isMuted);
  const handleVolumeChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseFloat(event.target.value);
    setVolume(value);
    if (value > 0) setIsMuted(false);
  };
  const handleSeek = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseFloat(event.target.value);
    seek(value);
    setProgress(value);
  };
  const formatTime = (seconds: number) => {
    if (!Number.isFinite(seconds) || seconds <= 0) return '0:00';
    return `${Math.floor(seconds / 60)}:${Math.floor(seconds % 60).toString().padStart(2, '0')}`;
  };
  if (!tracks.length) return null;
  return (
    <div className="bg-gray-50 border border-gray-200 rounded-lg p-3 sm:p-4 mb-4">
      <div className="flex items-center justify-between mb-2">
        <span className="text-xs sm:text-sm font-medium text-gray-700">
          {activeTrack?.label || verseLabel} - {currentTrack === 'arabic' ? 'Arabic' : currentTrack === 'english' ? 'English' : 'Audio'}
        </span>
        <span className="text-xs text-gray-500">
          {formatTime(progress)} / {formatTime(duration)}
        </span>
      </div>

      {/* Progress bar */}
      <input
        type="range"
        min="0"
        max={duration || 100}
        value={progress}
        onChange={handleSeek}
        disabled={!currentTrack}
        className="w-full h-1 mb-3 bg-gray-300 rounded-lg appearance-none cursor-pointer accent-black disabled:opacity-50"
        aria-label="Seek audio position"
      />

      <div className="flex items-center justify-between gap-2">
        {/* Playback controls */}
        <div className="flex items-center gap-1 sm:gap-2">
          {tracks.some(track => track.label === (activeTrack?.label || tracks[0]?.label) && track.language === 'arabic') && (
            <button
              onClick={skipToArabic}
              disabled={currentTrack === 'arabic' && isPlaying}
              className="p-1.5 sm:p-2 hover:bg-gray-200 rounded-full transition-colors disabled:opacity-50"
              title="Play Arabic"
              aria-label="Play Arabic recitation"
            >
              <SkipBack className="size-4 sm:size-5" />
            </button>
          )}
          
          <button
            onClick={togglePlay}
            className="p-2 sm:p-3 bg-black text-white hover:bg-gray-800 rounded-full transition-colors"
            title={isPlaying ? 'Pause' : 'Play'}
            aria-label={isPlaying ? 'Pause audio' : 'Play audio'}
          >
            {isPlaying ? (
              <Pause className="size-4 sm:size-5" />
            ) : (
              <Play className="size-4 sm:size-5" />
            )}
          </button>

          {tracks.some(track => track.label === (activeTrack?.label || tracks[0]?.label) && track.language === 'english') && (
            <button
              onClick={skipToEnglish}
              disabled={currentTrack === 'english' && isPlaying}
              className="p-1.5 sm:p-2 hover:bg-gray-200 rounded-full transition-colors disabled:opacity-50"
              title="Play English"
              aria-label="Play English translation"
            >
              <SkipForward className="size-4 sm:size-5" />
            </button>
          )}
        </div>

        {/* Volume controls */}
        <div className="flex items-center gap-1 sm:gap-2">
          <button
            onClick={toggleMute}
            className="p-1.5 sm:p-2 hover:bg-gray-200 rounded-full transition-colors"
            title={isMuted ? 'Unmute' : 'Mute'}
            aria-label={isMuted ? 'Unmute audio' : 'Mute audio'}
          >
            {isMuted ? (
              <VolumeX className="size-4 sm:size-5" />
            ) : (
              <Volume2 className="size-4 sm:size-5" />
            )}
          </button>
          <input
            type="range"
            min="0"
            max="1"
            step="0.1"
            value={volume}
            onChange={handleVolumeChange}
            className="w-16 sm:w-20 h-1 bg-gray-300 rounded-lg appearance-none cursor-pointer accent-black"
            title="Volume"
            aria-label="Volume control"
          />
        </div>
      </div>

      <audio
        ref={audioRef}
        preload="auto"
        onPlay={() => playbackChanged(true)}
        onPause={() => playbackChanged(false)}
        onTimeUpdate={() => { if (!streamRef.current) setProgress(audioRef.current?.currentTime || 0); }}
        onLoadedMetadata={() => { if (!streamRef.current) setDuration(Number.isFinite(audioRef.current?.duration) ? audioRef.current!.duration : 0); }}
        onEnded={handleEnded}
        onError={() => { playbackChanged(false); setError('Audio could not be loaded. Press play to retry.'); }}
      />
      <AudioDiagnostics audioRef={audioRef} />
      {error && <p role="status" className="text-xs text-gray-600 mt-2">{error}</p>}
    </div>
  );
}
