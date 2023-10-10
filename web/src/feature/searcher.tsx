import { Search } from "@mui/icons-material";
import { IconButton, Input } from "@mui/joy";

export function Searcher() {
  return (
    <Input
      placeholder="Search"
      endDecorator={
        <IconButton size="sm">
          <Search />
        </IconButton>
      }
      size="sm"
      style={{ width: "25%" }}
      sx={{
        ".MuiInput-input": {
          textAlign: "center",
        },
      }}
    />
  );
}
