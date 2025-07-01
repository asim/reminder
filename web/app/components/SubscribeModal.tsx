import React from 'react';

interface Props {
  open: boolean;
  onClose: () => void;
}

export function SubscribeModal({ open, onClose }: Props) {
  if (!open) return null;

  return (
    <div className='fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-40'>
      <div className='bg-white rounded-lg shadow-lg p-6 w-full max-w-sm'>
        <button className='absolute top-2 right-3 text-gray-400 hover:text-black' onClick={onClose}>&times;</button>
        <h2 className='text-xl font-semibold mb-3'>Subscribe to Daily Reminder</h2>
        <p className='text-gray-700'>We are migrating to web push notifications. Please ensure you have enabled notifications in your browser.</p>
      </div>
    </div>
  );
}
