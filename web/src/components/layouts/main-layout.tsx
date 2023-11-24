import React from "react";

import { Box, Grid, GridItem } from "@chakra-ui/react";

type MainLayoutProps = {
  appbar?: React.ReactNode;
  sidebar?: React.ReactNode;
  children: React.ReactNode;
};

const MainLayout = ({ appbar, sidebar, children }: MainLayoutProps) => {
  if (appbar === undefined && sidebar === undefined) {
    return <Box>{children}</Box>;
  }

  return (
    <Grid
      h="100vh"
      templateRows="repeat(12, 1fr)"
      templateColumns="repeat(12, 1fr)"
      gap={2}
    >
      {appbar && (
        <GridItem rowSpan={1} colSpan={12}>
          {appbar}
        </GridItem>
      )}

      {sidebar && (
        <GridItem rowSpan={appbar ? 11 : 12} colSpan={2}>
          {sidebar}
        </GridItem>
      )}

      <GridItem rowSpan={appbar ? 11 : 12} colSpan={sidebar ? 10 : 12}>
        {children}
      </GridItem>
    </Grid>
  );
};

export default MainLayout;
