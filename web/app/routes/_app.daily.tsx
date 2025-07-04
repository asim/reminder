import { useQuery } from '@tanstack/react-query';
import { httpGet, httpPost } from '~/utils/http';
import React, { useState, useEffect } from 'react';
import { subscribeUserToPush } from '../utils/push';
import { unsubscribeUserFromPush } from '../utils/push-unsub';
import { Outlet } from 'react-router';
import DailySidebarNav from '../components/daily-sidebar-nav';

interface DailyResponse {
  name: string;
  hadith: string;
  verse: string;
  links: Record<string, string>;
  updated: string;
  message: string;
}

export default function DailyLayout() {
  // Match Quran/Hadith layout: sidebar and content in flex-row, content scrollable, padding
  return (
    <div className="flex flex-row h-full">
      <DailySidebarNav />
      <div className="flex flex-col overflow-y-auto flex-1 px-5 py-5">
        <Outlet />
      </div>
    </div>
  );
}
