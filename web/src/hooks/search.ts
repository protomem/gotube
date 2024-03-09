import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";

const useSearch = (value?: string) => {
  const nav = useNavigate();
  const { handleSubmit: formHandleSubmit, register } = useForm();

  const inputProps = {
    id: "term",
    type: "text",
    autoComplete: "off",
    ...register("term", {
      value,
      required: "Please enter a search term",
      minLength: {
        value: 3,
        message: "Search term must be at least 3 characters",
      },
    }),
  };

  const handleSubmit = formHandleSubmit(({ term }) => {
    nav(`/search?term=${term}`);
  });

  return {
    handleSubmit,
    inputProps,
  };
};

export default useSearch;
