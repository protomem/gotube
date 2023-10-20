import { Button, ButtonGroup } from "@mui/joy";

export interface PaginatorProps {
  page: number;
  onChangePage: (page: number) => void;
  maxPage: number;
}

export default function Pagintor({
  page,
  onChangePage,
  maxPage,
}: PaginatorProps) {
  return (
    <ButtonGroup>
      {maxPage > 1 &&
        Array.from({ length: maxPage }, (_, i) => i + 1).map((i) => (
          <Button
            key={i}
            variant={page === i ? "solid" : "outlined"}
            color={page === i ? "primary" : "neutral"}
            onClick={() => onChangePage(i)}
          >
            {i}
          </Button>
        ))}
    </ButtonGroup>
  );
}
