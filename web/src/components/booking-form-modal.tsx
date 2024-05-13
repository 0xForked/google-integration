import {useEffect, useState} from "react";
import {
    AlertDialog, AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogHeader,
    AlertDialogTitle
} from "@/components/ui/alert-dialog.tsx";
import {Calendar, Clock12, Globe, Loader2, Video, X} from "lucide-react";
import {Avatar, AvatarFallback} from "@/components/ui/avatar.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Input} from "@/components/ui/input.tsx";
import {useForm} from "react-hook-form";
import {bookingSchema,TBookingSchema} from "@/lib/schema.ts";
import {yupResolver} from "@hookform/resolvers/yup";
import {Button} from "@/components/ui/button.tsx";
import {newBooking} from "@/lib/api.ts";
import {useNavigate} from "react-router-dom";
import {addMinutes, stringTimeToInt} from "@/lib/time.ts";

interface BookingFormDialogProps {
    display: boolean
    username?: string | null
    eventType?: EventType | null
    date?: Date | undefined
    time?: string | undefined
    callback(): void
}

export function BookingFormModal(props: BookingFormDialogProps) {
    const [display, setDisplay] = useState(false)
    const navigate = useNavigate()

    const {
        register,
        handleSubmit,
        formState: {errors, isSubmitting},
    } = useForm<TBookingSchema>({
        resolver: yupResolver(bookingSchema),
    })

    useEffect(() => {
        if (props.display && !props.eventType) {
            alert("event type data is required")
            props.callback()
            return
        }
        setDisplay(props.display)
    }, [props])

    const onSubmit = (data: TBookingSchema) => {
        newBooking(
            props.username!, props.eventType!.id,
            (Math.floor(props.date!.getTime() / 1000)),
            stringTimeToInt(props.time!), data.name,
            data.email, data.notes,
        ).then((resp) => {
            console.log(resp)
            navigate("/schedule/1")
        }).catch((err) => console.log(err))
    }

    return (<AlertDialog open={display} onOpenChange={props.callback}>
        <AlertDialogContent className="min-w-[700px]">
            <AlertDialogHeader className="flex flex-row justify-between">
                <section>
                    <AlertDialogTitle>Booking</AlertDialogTitle>
                    <AlertDialogDescription>
                        Please fill all the form.
                    </AlertDialogDescription>
                </section>
                <AlertDialogCancel className="border-none">
                    <X className="w-[21px]"/>
                </AlertDialogCancel>
            </AlertDialogHeader>
            <hr className="mb-4"/>
            <section className="grid grid-cols-3 divide-x-[1px]">
                <div className="px-2 flex flex-col gap-2">
                    <Avatar className="mb-2 h-8 w-8 text-xs">
                        <AvatarFallback>
                            {props?.username?.substring(0, 2).toUpperCase() ?? "-"}
                        </AvatarFallback>
                    </Avatar>
                    <h3 className="text-lg font-bold">{props?.eventType?.title}</h3>
                    <p className="text-sm font-light flex items-center">
                        <Clock12 className="w-[12px] h-[12px] mr-2"/>
                        {props?.eventType?.duration} mins
                    </p>
                    <p className="text-sm font-light flex items-center">
                        <Video className="w-[12px] h-[12px] mr-2"/>
                        Google Meet
                    </p>
                    <p className="text-sm font-light flex items-center">
                        <Globe className="w-[12px] h-[12px] mr-2"/>
                        {props?.eventType?.availability?.timezone}
                    </p>
                    <p className="text-sm font-light flex items-baseline">
                        <Calendar className="w-[12px] h-[12px] mr-2"/>
                        <span>
                            {props?.date?.toLocaleString('en-US',{ weekday: 'long' })},
                            <span className="ml-1">{props?.date?.toLocaleString('en-US', {month: "short"})}</span>
                            <span className="ml-1">{props?.date?.getDate()}</span>
                            <span className="ml-1">{props?.date?.getFullYear()}</span>
                            <span className="block">
                                {props?.time} -
                                {(props?.time && props.eventType) && addMinutes(
                                    props.time, props.eventType.duration
                                )}
                            </span>
                        </span>
                    </p>
                </div>
                <div className="col-span-2 min-h-80 px-6">
                    <form className="flex flex-col w-full" onSubmit={handleSubmit(onSubmit)}>
                        <div className="space-y-2 py-2 text-left">
                            <Label htmlFor="name">Your Name *</Label>
                            <Input
                                id="name"
                                placeholder="e.g: lorem ipsum"
                                {...register('name')}
                            />
                            <p className={`text-sm text-muted-foreground ${errors?.name ? "text-red-500" : ""}`}>
                                {errors?.name ? errors?.name?.message : "Enter your name"}
                            </p>
                        </div>
                        <div className="space-y-2 py-2 text-left">
                            <Label htmlFor="email">Email Address *</Label>
                            <Input
                                id="email"
                                placeholder="e.g: lorem@ipsum.id"
                                {...register('email')}
                            />
                            <p className={`text-sm text-muted-foreground ${errors?.email ? "text-red-500" : ""}`}>
                                {errors?.email ? errors?.email?.message : "Enter your email address"}
                            </p>
                        </div>
                        <div className="space-y-2 py-2 text-left">
                            <div className="space-y-2">
                                <Label htmlFor="notes">Additional notes</Label>
                                <Input
                                    id="notes"
                                    placeholder="Please share anything that will help prepare for our meeting."
                                    {...register('notes')}
                                />
                            </div>
                        </div>

                        <Button className="mt-4 ml-auto" type="submit" disabled={isSubmitting}>
                            {isSubmitting ? <Loader2 className="mr-2 h-4 w-4 animate-spin"/> : <></>}
                            Confirm
                        </Button>
                    </form>
                </div>
            </section>
        </AlertDialogContent>
        </AlertDialog>
    )
}