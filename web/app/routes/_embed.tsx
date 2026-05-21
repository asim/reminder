import { Outlet } from 'react-router';

export default function EmbedLayout() {
  return (
    <div className='embed-root'>
      <Outlet />
    </div>
  );
}
