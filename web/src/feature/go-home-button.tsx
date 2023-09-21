import React from "react";
import { ArrowBackIos } from "@mui/icons-material";
import { Button } from "@mui/joy";
import { useNavigate } from "react-router-dom";

export default function GoHomeButton() {
  const nav = useNavigate();

  const handleClick = (e: React.MouseEvent) => {
    e.preventDefault();
    nav("/", { replace: true });
  };

  return (
    <Button onClick={handleClick}>
      <ArrowBackIos /> Go Home
    </Button>
  );
}
