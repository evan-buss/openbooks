import {
  Button,
  Center,
  Paper,
  PasswordInput,
  Text,
  TextInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useState } from "react";
import useAuth from "../util/auth";

export function LogInPage() {
  const [error, setError] = useState<string | undefined>(undefined);
  const [loading, setLoading] = useState(false);
  const { logIn } = useAuth();

  const form = useForm({
    initialValues: {
      username: "",
      password: "",
    },
  });

  type FormValues = typeof form.values;

  const onSubmit = async (values: FormValues) => {
    try {
      setLoading(true);
      await logIn(values);
    } catch (err) {
      setLoading(false);
      setError(err as string);
    }
  };

  return (
    <Center sx={{ height: "100%" }}>
      <Paper sx={{ minWidth: "400px" }} withBorder p="xl">
        <form onSubmit={form.onSubmit(onSubmit)}>
          <TextInput
            required
            label="Username"
            placeholder="admin"
            {...form.getInputProps("username")}
          />
          <PasswordInput
            required
            mt="md"
            label="Password"
            type="password"
            placeholder="Password"
            {...form.getInputProps("password")}
          />
          {error && <Text color="red">{error}</Text>}
          <Button type="submit" mt="md" variant="filled" loading={loading}>
            Log In
          </Button>
        </form>
      </Paper>
    </Center>
  );
}
