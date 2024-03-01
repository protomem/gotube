import { ChangeEvent } from "react";
import {
  Box,
  Button,
  Flex,
  FormControl,
  FormLabel,
  Grid,
  Input,
  SimpleGrid,
  Switch,
  Textarea,
  useToast,
} from "@chakra-ui/react";
import { Field, Formik } from "formik";

type FormValues = {
  title: string;
  description?: string;
  thumbnailPath: string;
  videoPath: string;
  isPublic?: boolean;
};

const initialValues: FormValues = {
  title: "",
  description: "",
  thumbnailPath: "",
  videoPath: "",
  isPublic: true,
};

type Props = {
  afterClose?: () => void;
};

export default function UploadVideoForm({ afterClose }: Props) {
  const toast = useToast();

  const handleSubmit = (values: FormValues) => {
    console.log(values);

    afterClose?.();
  };

  return (
    <Box display="flex" justifyContent="center" alignItems="center">
      <Formik initialValues={initialValues} onSubmit={handleSubmit}>
        {({ handleSubmit, values, setValues }) => (
          <form onSubmit={handleSubmit}>
            <Flex gap="12" direction="row">
              <Flex gap="6" direction="column" minW="25rem">
                <FormControl>
                  <Field
                    as={Input}
                    id="title"
                    name="title"
                    placeholder="Title"
                    variant="filled"
                  />
                </FormControl>

                <FormControl>
                  <Field
                    as={Textarea}
                    id="description"
                    name="description"
                    placeholder="Description ..."
                    variant="filled"
                  />
                </FormControl>

                <FormControl>
                  <FormLabel htmlFor="isPublic">Public</FormLabel>
                  <Field
                    as={Switch}
                    id="isPublic"
                    name="isPublic"
                    size="lg"
                    onChange={(e: ChangeEvent<HTMLInputElement>) => {
                      setValues({ ...values, isPublic: e.target.checked });
                    }}
                    isChecked={values.isPublic}
                  />
                </FormControl>
              </Flex>

              <Flex gap="10" direction="column">
                <FormControl>
                  <FormLabel htmlFor="thumbnailPath">Thumbnail</FormLabel>
                  <input
                    type="file"
                    onChange={(e) => {
                      setValues({
                        ...values,
                        thumbnailPath: e.target.files?.[0].name ?? "",
                      });
                    }}
                  />
                </FormControl>

                <FormControl>
                  <FormLabel htmlFor="videoPath">Video</FormLabel>
                  <input
                    type="file"
                    onChange={(e) => {
                      setValues({
                        ...values,
                        videoPath: e.target.files?.[0].name ?? "",
                      });
                    }}
                  />
                </FormControl>
              </Flex>
            </Flex>

            <Button type="submit" mt="6">
              Save
            </Button>
          </form>
        )}
      </Formik>
    </Box>
  );
}
