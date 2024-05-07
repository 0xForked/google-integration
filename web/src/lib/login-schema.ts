import * as yup from "yup";

const usernameSpec = yup.string()
    .required("Username is required");

const passwordSpec = yup.string()
    .required("Password is required");

const loginSchema = yup.object({
    username: usernameSpec,
    password: passwordSpec
}).required();

type TLoginSchema = yup.InferType<typeof loginSchema>

export {
    loginSchema,
}

export type {
    TLoginSchema,
}