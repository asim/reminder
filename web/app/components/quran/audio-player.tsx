import { Pause, Play, SkipBack, SkipForward, Volume2, VolumeX } from 'lucide-react';
import React, { useEffect, useRef, useState } from 'react';

type AudioPlayerProps = {
  arabicUrl?: string;
  englishUrl?: string;
  verseLabel: string;
  autoPlay?: boolean;
  onPlayComplete?: () => void;
};

export function AudioPlayer({
  arabicUrl,
  englishUrl,
  verseLabel,
  autoPlay = false,
  onPlayComplete,
}: AudioPlayerProps) {
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTrack, setCurrentTrack] = useState<'arabic' | 'english' | null>(null);
  const [progress, setProgress] = useState(0);
  const [duration, setDuration] = useState(0);
  const [volume, setVolume] = useState(1);
  const [isMuted, setIsMuted] = useState(false);
  const audioRef = useRef<HTMLAudioElement>(null);

  // Reset player when URLs change
  useEffect(() => {
    setIsPlaying(false);
    setCurrentTrack(null);
    setProgress(0);
  }, [arabicUrl, englishUrl]);

  // Handle playback sequence: Arabic first, then English
  const playSequence = async () => {
    if (!arabicUrl && !englishUrl) return;

    setIsPlaying(true);
    
    // Start with Arabic if available
    if (arabicUrl) {
      setCurrentTrack('arabic');
    } else if (englishUrl) {
      setCurrentTrack('english');
    }
  };

  const pause = () => {
    if (audioRef.current) {
      audioRef.current.pause();
    }
    setIsPlaying(false);
  };

  const togglePlay = () => {
    if (isPlaying) {
      pause();
    } else {
      playSequence();
    }
  };

  const skipToEnglish = () => {
    if (englishUrl) {
      if (audioRef.current) {
        audioRef.current.pause();
      }
      setCurrentTrack('english');
      setIsPlaying(true);
    }
  };

  const skipToArabic = () => {
    if (arabicUrl) {
      if (audioRef.current) {
        audioRef.current.pause();
      }
      setCurrentTrack('arabic');
      setIsPlaying(true);
    }
  };

  const toggleMute = () => {
    if (audioRef.current) {
      audioRef.current.muted = !isMuted;
      setIsMuted(!isMuted);
    }
  };

  const handleVolumeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newVolume = parseFloat(e.target.value);
    setVolume(newVolume);
    if (audioRef.current) {
      audioRef.current.volume = newVolume;
    }
    if (newVolume > 0 && isMuted) {
      setIsMuted(false);
      if (audioRef.current) {
        audioRef.current.muted = false;
      }
    }
  };

  const handleTimeUpdate = () => {
    if (audioRef.current) {
      setProgress(audioRef.current.currentTime);
    }
  };

  const handleLoadedMetadata = () => {
    if (audioRef.current) {
      setDuration(audioRef.current.duration);
    }
  };

  const handleEnded = () => {
    // When Arabic ends, play English
    if (currentTrack === 'arabic' && englishUrl) {
      setCurrentTrack('english');
    } else {
      // Sequence complete
      setIsPlaying(false);
      setCurrentTrack(null);
      setProgress(0);
      if (onPlayComplete) {
        onPlayComplete();
      }
    }
  };

  const handleSeek = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newTime = parseFloat(e.target.value);
    setProgress(newTime);
    if (audioRef.current) {
      audioRef.current.currentTime = newTime;
    }
  };

  const formatTime = (seconds: number) => {
    if (isNaN(seconds)) return '0:00';
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const currentUrl = currentTrack === 'arabic' ? arabicUrl : currentTrack === 'english' ? englishUrl : null;

  if (!arabicUrl && !englishUrl) {
    return null;
  }

  return (
    <div className="bg-gray-50 border border-gray-200 rounded-lg p-3 sm:p-4 mb-4">
      <div className="flex items-center justify-between mb-2">
        <span className="text-xs sm:text-sm font-medium text-gray-700">
          {verseLabel} - {currentTrack === 'arabic' ? 'Arabic' : currentTrack === 'english' ? 'English' : 'Audio'}
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
        className="w-full h-1 mb-3 bg-gray-300 rounded-lg appearance-none cursor-pointer accent-black"
      />

      <div className="flex items-center justify-between gap-2">
        {/* Playback controls */}
        <div className="flex items-center gap-1 sm:gap-2">
          {arabicUrl && (
            <button
              onClick={skipToArabic}
              disabled={currentTrack === 'arabic' && isPlaying}
              className="p-1.5 sm:p-2 hover:bg-gray-200 rounded-full transition-colors disabled:opacity-50"
              title="Play Arabic"
            >
              <SkipBack className="size-4 sm:size-5" />
            </button>
          )}
          
          <button
            onClick={togglePlay}
            className="p-2 sm:p-3 bg-black text-white hover:bg-gray-800 rounded-full transition-colors"
            title={isPlaying ? 'Pause' : 'Play'}
          >
            {isPlaying ? (
              <Pause className="size-4 sm:size-5" />
            ) : (
              <Play className="size-4 sm:size-5" />
            )}
          </button>

          {englishUrl && (
            <button
              onClick={skipToEnglish}
              disabled={currentTrack === 'english' && isPlaying}
              className="p-1.5 sm:p-2 hover:bg-gray-200 rounded-full transition-colors disabled:opacity-50"
              title="Play English"
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
          />
        </div>
      </div>

      {/* Hidden audio element */}
      {currentUrl && (
        <audio
          ref={audioRef}
          src={currentUrl}
          onTimeUpdate={handleTimeUpdate}
          onLoadedMetadata={handleLoadedMetadata}
          onEnded={handleEnded}
          autoPlay={isPlaying}
          volume={volume}
          muted={isMuted}
        />
      )}
    </div>
  );
}
