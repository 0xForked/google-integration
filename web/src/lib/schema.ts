import * as yup from "yup";

const usernameSpec = yup.string()
    .required("Username is required");

const passwordSpec = yup.string()
    .required("Password is required");

const nameSpec = yup.string()
    .required("Name is required");

const emailSpec = yup.string()
    .required("Email is required");

const loginSchema = yup.object({
    username: usernameSpec,
    password: passwordSpec
}).required();

const bookingSchema = yup.object({
    name: nameSpec,
    email: emailSpec,
    notes: yup.string()
}).required();

type TLoginSchema = yup.InferType<typeof loginSchema>
type TBookingSchema = yup.InferType<typeof bookingSchema>

export {
    loginSchema,
    bookingSchema
}

export type {
    TLoginSchema,
    TBookingSchema
}