import {useEffect, useState} from "react";
import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogHeader,
    AlertDialogTitle
} from "@/components/ui/alert-dialog.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Loader2} from "lucide-react";
import {loginSchema, TLoginSchema} from "@/lib/login-schema.ts";
import {useForm} from "react-hook-form";
import {yupResolver} from '@hookform/resolvers/yup';
import {Button} from "@/components/ui/button.tsx";
import {signIn} from "@/lib/api.ts";

interface LoginDialogProps {
    display: boolean,
}

export function LoginDialog(props: LoginDialogProps) {
    const [display, setDisplay] = useState(false)

    useEffect(() => setDisplay(props.display), [props])

    const {
        register,
        handleSubmit,
        formState: {errors, isSubmitting},
        setError
    } = useForm<TLoginSchema>({
        resolver: yupResolver(loginSchema),
    })

    const onSubmit = (data: TLoginSchema) => {
        signIn(data.username, data.password)
            .then((resp) => {
                if (resp.error && resp.error.includes("username")) {
                    setError("username", {
                        message: resp.error
                    })
                    return
                }
                if (resp.error && resp.error.includes("password")) {
                    setError("password", {
                        message: resp.error
                    })
                    return
                }
                if (resp.data.token != "") {
                    window.location.reload()
                }
            })
            .catch((err) => console.log(err))
    }

    return <>
        <AlertDialog open={display}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Login</AlertDialogTitle>
                    <AlertDialogDescription>
                        Sign in to your account to continue
                    </AlertDialogDescription>
                    <form onSubmit={handleSubmit(onSubmit)}>
                        <div className="space-y-4 py-2 text-left">
                            <div className="space-y-2">
                                <Label htmlFor="username">username</Label>
                                <Input
                                    id="username"
                                    placeholder="e.g: lorem_ipsum"
                                    {...register('username')}
                                />
                                <p className={`text-sm text-muted-foreground ${errors?.username ? "text-red-500" : ""}`}>
                                    {errors?.username ? errors?.username?.message : "Enter your username address"}
                                </p>
                            </div>
                        </div>
                        <div className="space-y-4 py-2 pb-4 text-left">
                            <div className="space-y-2">
                                <Label htmlFor="password">Password</Label>
                                <Input
                                    id="password"
                                    type="password"
                                    placeholder="* * * * *"
                                    {...register('password')}
                                />
                                <p className={`text-sm text-muted-foreground ${errors?.password ? "text-red-500" : ""}`}>
                                    {errors?.password ? errors?.password?.message : "Enter your password"}
                                </p>
                            </div>
                        </div>
                        <Button type="submit" disabled={isSubmitting}>
                            {isSubmitting ? <Loader2 className="mr-2 h-4 w-4 animate-spin"/> : <></>}
                            Sign In
                        </Button>
                    </form>
                </AlertDialogHeader>
            </AlertDialogContent>
        </AlertDialog>
    </>
}