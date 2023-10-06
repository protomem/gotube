import React from "react";
import { Box } from "@mui/joy";

export interface LayoutProps {
  children: React.ReactNode[];
  withSideBar?: boolean;
}

export default function Layout({ children, withSideBar }: LayoutProps) {
  if (withSideBar === undefined) withSideBar = true;

  return (
    <Box
      style={{
        width: "100vw",
        height: "100vh",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "start",
      }}
    >
      {children.length > 0 && (
        <Box
          style={{
            width: "100%",
            height: "7%",
            borderBottom: "1px solid gray",
            display: "flex",
            flexDirection: "row",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          {children[0]}
        </Box>
      )}

      {children.length > 1 && (
        <Box
          style={{
            width: "100%",
            height: "93%",
            display: "flex",
            flexDirection: "row",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          {!withSideBar ? (
            children[1]
          ) : (
            <>
              <Box
                style={{
                  width: "20%",
                  height: "100%",
                  overflowY: "auto",
                  borderRight: "1px solid gray",
                }}
              >
                {children[1]}
              </Box>

              <Box
                style={{
                  width: "80%",
                  height: "100%",
                  display: "flex",
                  flexDirection: "column",
                  overflowY: "auto",
                }}
              >
                {children.length > 2 &&
                  children
                    .slice(2)
                    .map((child, index) => <Box key={index}>{child}</Box>)}
              </Box>
            </>
          )}
        </Box>
      )}
    </Box>
  );
}
