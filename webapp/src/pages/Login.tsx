import { Button } from "@/components/ui/button";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { useAuth } from "@/hooks/useAuth";
import { useNavigate } from "react-router-dom";
import * as z from "zod"
import { useForm } from "@tanstack/react-form"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import AuthLayout from "@/components/layouts/AuthLayout";
import { normalizeErrors } from "@/lib/utils";

const formSchema = z.object({
    name: z
        .string()
        .min(1, "Введите логин"),
    password: z
        .string()
        .min(1, "Введите пароль"),
});

const Login: React.FC = () => {
    const navigate = useNavigate();
    const { login } = useAuth();

    const form = useForm({
        defaultValues: {
            name: "",
            password: "",
        },
        validators: {
            onSubmit: formSchema,
        },
        onSubmit: async ({ value }) => {
            try {
                await login(value.name, value.password);
                navigate("/admin/dashboard");
            } catch (error) {
                const errMsg = (error as Error).message
                if (errMsg.toLowerCase().includes("invalid password")) {
                    form.setFieldMeta("password", (meta) => ({
                        ...meta,
                        errorMap: { onSubmit: "Неверный пароль" }
                    }));
                } else if (errMsg.toLowerCase().includes("user not found")) {
                    form.setFieldMeta("name", (meta) => ({
                        ...meta,
                        errorMap: { onSubmit: "Данного пользователя не существует" }
                    }));
                }
            }
        },
    });

    return (
        <AuthLayout>
            <Card className="w-full sm:max-w-md">
                <CardHeader>
                    <CardTitle>Оконно</CardTitle>
                    <CardDescription>
                        Добро пожаловать в админ панель.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <form
                        id="dashboard-login-form"
                        onSubmit={(e) => {
                            e.preventDefault();
                            form.handleSubmit();
                        }}
                    >
                        <FieldGroup>
                            <form.Field
                                name="name"
                                children={(field) => {
                                    const error = field.state.meta.errorMap.onChange || field.state.meta.errorMap.onSubmit
                                    const isInvalid = Boolean(error);
                                    return (
                                        <Field data-invalid={isInvalid}>
                                            <FieldLabel htmlFor={field.name}>Логин</FieldLabel>
                                            <Input
                                                id={field.name}
                                                name={field.name}
                                                value={field.state.value}
                                                onBlur={field.handleBlur}
                                                onChange={(e) => field.handleChange(e.target.value)}
                                                data-invalid={isInvalid}
                                                placeholder="admin"
                                                autoComplete="off"
                                                required
                                            />
                                            {error && (
                                                <FieldError errors={normalizeErrors(error)} />
                                            )}
                                        </Field>
                                    )
                                }}
                            />
                            <form.Field
                                name="password"
                                children={(field) => {
                                    const error = field.state.meta.errorMap.onChange || field.state.meta.errorMap.onSubmit
                                    const isInvalid = Boolean(error);
                                    return (
                                        <Field data-invalid={isInvalid}>
                                            <FieldLabel htmlFor={field.name}>Пароль</FieldLabel>
                                            <Input
                                                id={field.name}
                                                name={field.name}
                                                value={field.state.value}
                                                onBlur={field.handleBlur}
                                                onChange={(e) => field.handleChange(e.target.value)}
                                                data-invalid={isInvalid}
                                                placeholder="1234"
                                                autoComplete="off"
                                                type="password"
                                                required
                                            />
                                            {error && (
                                                <FieldError errors={normalizeErrors(error)} />
                                            )}
                                        </Field>
                                    )
                                }}
                            />
                        </FieldGroup>
                    </form>
                </CardContent>
                <CardFooter>
                    <Field orientation="horizontal">
                        <Button
                            type="submit"
                            form="dashboard-login-form"
                            disabled={form.state.isFormValidating || form.state.isSubmitting}
                        >
                            Войти
                            {(form.state.isFormValidating || form.state.isSubmitting) && <Spinner />}
                        </Button>
                    </Field>
                </CardFooter>
            </Card>
        </AuthLayout>
    );
};

export default Login;
