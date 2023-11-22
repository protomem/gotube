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

  // TODO: add handle missing appbar or sidebar
  return (
    <Grid
      h="100dvh"
      templateRows="repeat(12, 1fr)"
      templateColumns="repeat(12, 1fr)"
      gap={2}
    >
      {appbar && (
        <GridItem rowSpan={1} colSpan={12} bg="tomato">
          {appbar}
        </GridItem>
      )}

      {sidebar && (
        <GridItem rowSpan={appbar ? 11 : 12} colSpan={2} bg="rosybrown">
          {sidebar}
        </GridItem>
      )}

      <GridItem
        rowSpan={appbar ? 11 : 12}
        colSpan={sidebar ? 10 : 12}
        bg="plum"
      >
        {children}
      </GridItem>
    </Grid>
  );
};

export default MainLayout;
