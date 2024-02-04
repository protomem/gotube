import { useFormik } from "formik";
import { FaSistrix } from "react-icons/fa6";
import {
  Box,
  IconButton,
  Input,
  InputGroup,
  InputRightAddon,
} from "@chakra-ui/react";
import { useRouter } from "next/navigation";

interface Props {
  defaultTerm?: string;
}

export default function Searcher({ defaultTerm = "" }: Props) {
  const router = useRouter()

  const formik = useFormik({
    initialValues: {
      term: defaultTerm,
    },
    onSubmit: ({ term }) => {
      if (term === "") return;
      router.push(`/search?term=${term}`);
    },
  });

  return (
    <Box w="2xl">
      <form onSubmit={formik.handleSubmit}>
        <InputGroup w="100%">
          <Input
            id="term"
            name="term"
            type="text"
            rounded="full"
            placeholder="Search ..."
            textAlign="center"
            onChange={formik.handleChange}
            value={formik.values.term}
            autoComplete="off"
          />
          <InputRightAddon rounded="full" px="2">
            <IconButton
              aria-label="Search"
              rounded="full"
              variant="ghost"
              icon={<FaSistrix />}
              type="submit"
            />
          </InputRightAddon>
        </InputGroup>
      </form>
    </Box>
  );
}
