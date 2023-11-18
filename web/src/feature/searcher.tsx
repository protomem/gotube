import { appendQuery } from "@/lib";
import { ROUTES } from "@/shared/constants/routes";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/router";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { Button } from "@/shared/ui/button";
import { Form, FormField } from "@/shared/ui/form";
import { Input } from "@/shared/ui/input";
import { MagnifyingGlassIcon } from "@radix-ui/react-icons";

const FormSchema = z.object({
  query: z.string().min(3, {
    message: "Please enter a search term",
  }),
});

export function Searcher() {
  const router = useRouter();

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      query: "",
    },
  });

  const onSubmit = (data: z.infer<typeof FormSchema>) => {
    router.push(appendQuery(ROUTES.SEARCH, { key: "q", value: data.query }));
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex flex-row justify-center items-center"
      >
        <FormField
          control={form.control}
          name="query"
          render={({ field }) => (
            <Input
              placeholder="Search..."
              className="h-8 rounded-xl text-center"
              {...field}
            />
          )}
        />
        <Button variant="ghost" size="icon" type="submit" className="h-7 w-10">
          <MagnifyingGlassIcon className="h-5 w-5" />
        </Button>
      </form>
    </Form>
  );
}
