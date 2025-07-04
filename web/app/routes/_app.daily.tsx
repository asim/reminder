import { useQuery } from '@tanstack/react-query';
import { httpGet, httpPost } from '~/utils/http';
import React, { useState, useEffect } from 'react';
import { subscribeUserToPush } from '../utils/push';
import { unsubscribeUserFromPush } from '../utils/push-unsub';
import { Outlet } from 'react-router';
import DailySidebar from './_app.daily.index';

interface DailyResponse {
  name: string;
  hadith: string;
  verse: string;
  links: Record<string, string>;
  updated: string;
  message: string;
}

export default function DailyLayout() {
  return (
    <div className="flex h-full">
      <div className="w-full lg:w-64">
        <DailySidebar />
      </div>
      <div className="flex-1">
        <Outlet />
      </div>
    </div>
  );
}
