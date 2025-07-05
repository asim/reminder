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
  // Defensive: render sidebar and content regardless of sidebar fetch errors
  return (
    <div className="flex flex-row h-full">
      <ErrorBoundary>
        <DailySidebarNav />
      </ErrorBoundary>
      <div className="flex flex-col flex-1 p-0 lg:p-8 mx-auto w-full lg:max-w-3xl space-y-8 overflow-y-auto">
        <Outlet />
      </div>
    </div>
  );
}

// Simple error boundary to prevent sidebar errors from breaking the page
class ErrorBoundary extends React.Component<{ children: React.ReactNode }, { hasError: boolean }> {
  constructor(props: { children: React.ReactNode }) {
    super(props);
    this.state = { hasError: false };
  }
  static getDerivedStateFromError() {
    return { hasError: true };
  }
  componentDidCatch(error: any, info: any) {
    // Optionally log error
  }
  render() {
    if (this.state.hasError) {
      return <div className="w-full lg:w-64 p-4 text-gray-400">Sidebar unavailable</div>;
    }
    return this.props.children;
  }
}
