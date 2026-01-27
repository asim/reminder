import { Pause, Play, SkipBack, SkipForward, Volume2, VolumeX } from 'lucide-react';
import React, { useEffect, useRef, useState } from 'react';

type AudioPlayerProps = {
  arabicUrl?: string;
  englishUrl?: string;
  verseLabel: string;
  onPlayComplete?: () => void;
  onPlayStart?: () => void;
  autoPlay?: boolean;
};

export function AudioPlayer({
  arabicUrl,
  englishUrl,
  verseLabel,
  onPlayComplete,
  onPlayStart,
  autoPlay = false,
}: AudioPlayerProps) {
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTrack, setCurrentTrack] = useState<'arabic' | 'english' | null>(null);
  const [progress, setProgress] = useState(0);
  const [duration, setDuration] = useState(0);
  const [volume, setVolume] = useState(1);
  const [isMuted, setIsMuted] = useState(false);
  const arabicAudioRef = useRef<HTMLAudioElement>(null);
  const englishAudioRef = useRef<HTMLAudioElement>(null);

  // Reset player when URLs change, and auto-play if enabled
  useEffect(() => {
    setCurrentTrack(null);
    setProgress(0);
    setDuration(0);
    
    if (autoPlay && (arabicUrl || englishUrl)) {
      // Small delay to ensure audio elements are ready
      const timer = setTimeout(() => {
        if (arabicUrl && arabicAudioRef.current) {
          setCurrentTrack('arabic');
          arabicAudioRef.current.currentTime = 0;
          arabicAudioRef.current.play().catch(console.error);
          setIsPlaying(true);
        } else if (englishUrl && englishAudioRef.current) {
          setCurrentTrack('english');
          englishAudioRef.current.currentTime = 0;
          englishAudioRef.current.play().catch(console.error);
          setIsPlaying(true);
        }
      }, 100);
      return () => clearTimeout(timer);
    } else {
      setIsPlaying(false);
    }
  }, [arabicUrl, englishUrl, autoPlay]);

  // Update audio element volume and muted state
  useEffect(() => {
    if (arabicAudioRef.current) {
      arabicAudioRef.current.volume = volume;
      arabicAudioRef.current.muted = isMuted;
    }
    if (englishAudioRef.current) {
      englishAudioRef.current.volume = volume;
      englishAudioRef.current.muted = isMuted;
    }
  }, [volume, isMuted]);

  // Update duration when track changes
  useEffect(() => {
    if (currentTrack === 'arabic' && arabicAudioRef.current) {
      const audio = arabicAudioRef.current;
      if (audio.readyState >= 1) {
        setDuration(audio.duration);
      }
    } else if (currentTrack === 'english' && englishAudioRef.current) {
      const audio = englishAudioRef.current;
      if (audio.readyState >= 1) {
        setDuration(audio.duration);
      }
    }
  }, [currentTrack]);

  // Handle playback sequence: Arabic first, then English
  const playSequence = async () => {
    if (!arabicUrl && !englishUrl) return;

    // Notify parent that playback started
    if (onPlayStart) {
      onPlayStart();
    }

    // Start with Arabic if available
    if (arabicUrl && arabicAudioRef.current) {
      setCurrentTrack('arabic');
      arabicAudioRef.current.currentTime = 0;
      setProgress(0);
      if (arabicAudioRef.current.readyState >= 1) {
        setDuration(arabicAudioRef.current.duration);
      }
      arabicAudioRef.current.play().catch(console.error);
      setIsPlaying(true);
    } else if (englishUrl && englishAudioRef.current) {
      setCurrentTrack('english');
      englishAudioRef.current.currentTime = 0;
      setProgress(0);
      if (englishAudioRef.current.readyState >= 1) {
        setDuration(englishAudioRef.current.duration);
      }
      englishAudioRef.current.play().catch(console.error);
      setIsPlaying(true);
    }
  };

  const pause = () => {
    if (arabicAudioRef.current) {
      arabicAudioRef.current.pause();
    }
    if (englishAudioRef.current) {
      englishAudioRef.current.pause();
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
    if (englishUrl && englishAudioRef.current) {
      // Pause Arabic
      if (arabicAudioRef.current) {
        arabicAudioRef.current.pause();
      }
      setCurrentTrack('english');
      setProgress(0);
      englishAudioRef.current.currentTime = 0;
      if (englishAudioRef.current.readyState >= 1) {
        setDuration(englishAudioRef.current.duration);
      }
      englishAudioRef.current.play().catch(console.error);
      setIsPlaying(true);
    }
  };

  const skipToArabic = () => {
    if (arabicUrl && arabicAudioRef.current) {
      // Pause English
      if (englishAudioRef.current) {
        englishAudioRef.current.pause();
      }
      setCurrentTrack('arabic');
      setProgress(0);
      arabicAudioRef.current.currentTime = 0;
      if (arabicAudioRef.current.readyState >= 1) {
        setDuration(arabicAudioRef.current.duration);
      }
      arabicAudioRef.current.play().catch(console.error);
      setIsPlaying(true);
    }
  };

  const toggleMute = () => {
    setIsMuted(!isMuted);
  };

  const handleVolumeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newVolume = parseFloat(e.target.value);
    setVolume(newVolume);
    if (newVolume > 0 && isMuted) {
      setIsMuted(false);
    }
  };

  const handleArabicTimeUpdate = () => {
    if (arabicAudioRef.current && currentTrack === 'arabic') {
      setProgress(arabicAudioRef.current.currentTime);
    }
  };

  const handleEnglishTimeUpdate = () => {
    if (englishAudioRef.current && currentTrack === 'english') {
      setProgress(englishAudioRef.current.currentTime);
    }
  };

  const handleArabicLoadedMetadata = () => {
    if (arabicAudioRef.current && currentTrack === 'arabic') {
      setDuration(arabicAudioRef.current.duration);
    }
  };

  const handleEnglishLoadedMetadata = () => {
    if (englishAudioRef.current && currentTrack === 'english') {
      setDuration(englishAudioRef.current.duration);
    }
  };

  const handleArabicEnded = () => {
    // When Arabic ends, play English
    if (englishUrl && englishAudioRef.current) {
      setCurrentTrack('english');
      setProgress(0);
      englishAudioRef.current.currentTime = 0;
      if (englishAudioRef.current.readyState >= 1) {
        setDuration(englishAudioRef.current.duration);
      }
      englishAudioRef.current.play().catch(console.error);
    } else {
      // No English, sequence complete
      setIsPlaying(false);
      setCurrentTrack(null);
      setProgress(0);
      setDuration(0);
      if (onPlayComplete) {
        onPlayComplete();
      }
    }
  };

  const handleEnglishEnded = () => {
    // English finished, sequence complete
    setIsPlaying(false);
    setCurrentTrack(null);
    setProgress(0);
    setDuration(0);
    if (onPlayComplete) {
      onPlayComplete();
    }
  };

  const handleSeek = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newTime = parseFloat(e.target.value);
    setProgress(newTime);
    if (currentTrack === 'arabic' && arabicAudioRef.current) {
      arabicAudioRef.current.currentTime = newTime;
    } else if (currentTrack === 'english' && englishAudioRef.current) {
      englishAudioRef.current.currentTime = newTime;
    }
  };

  const formatTime = (seconds: number) => {
    if (isNaN(seconds) || seconds === 0) return '0:00';
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

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
        disabled={!currentTrack}
        className="w-full h-1 mb-3 bg-gray-300 rounded-lg appearance-none cursor-pointer accent-black disabled:opacity-50"
        aria-label="Seek audio position"
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

          {englishUrl && (
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

      {/* Hidden audio elements - both always mounted for seamless switching */}
      {arabicUrl && (
        <audio
          ref={arabicAudioRef}
          src={arabicUrl}
          onTimeUpdate={handleArabicTimeUpdate}
          onLoadedMetadata={handleArabicLoadedMetadata}
          onEnded={handleArabicEnded}
          preload="auto"
        />
      )}
      {englishUrl && (
        <audio
          ref={englishAudioRef}
          src={englishUrl}
          onTimeUpdate={handleEnglishTimeUpdate}
          onLoadedMetadata={handleEnglishLoadedMetadata}
          onEnded={handleEnglishEnded}
          preload="auto"
        />
      )}
    </div>
  );
}
