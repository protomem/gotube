import React from "react";

interface MainLayoutProps {
  appbar: React.ReactNode;
  sidebar?: React.ReactNode;
  children: React.ReactNode;
}

export function MainLayout({ appbar, sidebar, children }: MainLayoutProps) {
  return (
    <div className="w-screen h-screen">
      <header className="w-auto h-[50px] flex justify-center border-b-[1px] border-b-neutral-800 dark:border-b-neutral-200">
        {appbar}
      </header>

      <div className="flex flex-row h-[calc(100%-50px)]">
        {sidebar && (
          <div className="w-72 h-full border-r-[1px] border-r-neutral-800 dark:border-r-neutral-200">
            {sidebar}
          </div>
        )}

        <main className="w-full h-full">{children}</main>
      </div>
    </div>
  );
}
