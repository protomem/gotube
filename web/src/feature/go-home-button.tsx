import React from "react";
import { ArrowBackIos } from "@mui/icons-material";
import { Button } from "@mui/joy";
import { useNavigate } from "react-router-dom";

export interface GoHomeButtonProps {
  withArrow?: boolean;
}

export default function GoHomeButton({ withArrow }: GoHomeButtonProps) {
  if (withArrow === undefined) withArrow = true;

  const nav = useNavigate();

  const handleClick = (e: React.MouseEvent) => {
    e.preventDefault();
    nav("/", { replace: true });
  };

  return (
    <Button onClick={handleClick}>
      {withArrow && <ArrowBackIos />} Go Home
    </Button>
  );
}
