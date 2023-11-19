import React from "react";

interface MainLayoutProps {
  appbar: React.ReactNode;
  sidebar?: React.ReactNode;
  children: React.ReactNode;
}

export function MainLayout({ appbar, sidebar, children }: MainLayoutProps) {
  return (
    <div className="w-screen h-screen">
      <header className="w-auto h-[50px] flex justify-center">{appbar}</header>

      <div className="flex flex-row h-[calc(100%-50px)]">
        {sidebar && <div className="w-[18rem] h-full">{sidebar}</div>}

        <main className="w-[calc(100%-18rem)] h-full">{children}</main>
      </div>
    </div>
  );
}
