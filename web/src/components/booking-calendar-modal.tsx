import {useEffect, useState} from "react";
import {
    AlertDialog, AlertDialogCancel,
    AlertDialogContent, AlertDialogDescription,
    AlertDialogHeader, AlertDialogTitle,
} from "@/components/ui/alert-dialog.tsx";
import {Avatar, AvatarFallback} from "@/components/ui/avatar.tsx";
import {Clock12, Globe, Video, X} from "lucide-react";
import {Calendar} from "@/components/ui/calendar.tsx";
import {Button} from "@/components/ui/button.tsx";
import {generateTimeRange, getInitialDate, isGreaterThanMaxEndTime, isSameDay} from "@/lib/time.ts";
import {RadioGroup, RadioGroupItem} from "@/components/ui/radio-group.tsx";
import {Label} from "@radix-ui/react-dropdown-menu";

interface BookingCalendarDialogProps {
    display: boolean
    username?: string | null
    eventType?: EventType | null
    callback(): void
    nextStep(date: Date, time: string, meetingSource: string): void
}

export function BookingCalendarModal(props: BookingCalendarDialogProps) {
    const [display, setDisplay] = useState(false)
    const [date, setDate] = useState<Date | undefined>(getInitialDate())
    const [blockDays, setBlockDays] = useState<number[]>([])
    const [startTime, setStartTime] = useState<number>(800)
    const [endTime, setEndTime] = useState<number>(1700)
    const [meetingSource, setMeetingSource] = useState("google")

    useEffect(() => {
        if (props.display && !props.eventType) {
            alert("event type data is required")
            props.callback()
            return
        }
        let days: number[] = []
        props.eventType?.availability?.days?.forEach((ad: AvailabilityDays) => {
            if (!ad.enable) {
                days.push(ad.day)
            }
        })
        setBlockDays(days)
        setSelectedDate(getInitialDate())
        setDisplay(props.display)
    }, [props])

    const setSelectedDate = (date: Date | undefined) => {
        if (date != undefined && blockDays?.includes(date.getDay())) {
            const newDate = new Date(date);
            if (date.getDay() == 0) {
                newDate.setDate(date.getDate() + 1);
            }
            if (date.getDay() == 6) {
                newDate.setDate(date.getDate() + 2);
            }
            date = newDate
        }
        setDate(date)
        props.eventType?.availability?.days?.forEach((ad) => {
            if (ad.day === date?.getDay()) {
                setStartTime(ad.start_time)
                setEndTime(ad.end_time)
            }
        })
    }

    const disabledCalendar = (date: Date) =>
        date.getTime() < new Date().setHours(0, 0, 0, 0) ||
            blockDays?.includes(date.getDay())

    const isTimeAvailable = (): boolean => {
        if (date == undefined) {
            return false
        }
        if (date.getTime() < new Date().setHours(0, 0, 0, 0)) {
            return false
        }
        if (!isSameDay(date!)) {
            return true
        }
        return !isGreaterThanMaxEndTime(date!, endTime)
    }

    const setNextStep = (time: string) =>  props.nextStep(date!, time, meetingSource)

    return (<AlertDialog open={display} onOpenChange={props.callback}>
        <AlertDialogContent className="min-w-[850px]">
            <AlertDialogHeader className="flex flex-row justify-between">
                <section>
                    <AlertDialogTitle>Booking</AlertDialogTitle>
                    <AlertDialogDescription>
                        Please select a date and time to proceed with the booking.
                    </AlertDialogDescription>
                </section>
                <AlertDialogCancel className="border-none">
                    <X className="w-[21px]"/>
                </AlertDialogCancel>
            </AlertDialogHeader>
            <hr className="mb-4"/>
            <section className="grid grid-cols-4 divide-x-[1px]">
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
                        Meeting Location
                    </p>
                    <RadioGroup
                      className="ml-6"
                      defaultValue={meetingSource}
                      onValueChange={(val: string) => setMeetingSource(val)}
                    >
                        {props?.eventType?.is_google_available &&<div className="flex items-center space-x-2">
                            <RadioGroupItem className="w-[12px] h-[12px]" value="google" />
                            <Label className="text-sm">Google Meet</Label>
                        </div>}
                        {props?.eventType?.is_microsoft_available && <div className="flex items-center space-x-2">
                            <RadioGroupItem className="w-[12px] h-[12px]" value="microsoft"/>
                            <Label className="text-sm">Microsoft Team</Label>
                        </div>}
                    </RadioGroup>
                            <p className="text-sm font-light flex items-center">
                                <Globe className="w-[12px] h-[12px] mr-2"/>
                                {props?.eventType?.availability?.timezone}
                            </p>
                        </div>
                        <div className="px-2 col-span-2 h-96 flex w-full">
                    <Calendar
                        mode="single"
                        selected={date}
                        onSelect={(date) => setSelectedDate(date)}
                        disabled={(date) => disabledCalendar(date)}
                        className="rounded-md mx-auto"
                    />
                </div>
                <div className="px-2">
                    <p className="font-bold inline">
                        {date?.toLocaleString('en-US',{ weekday: 'short' })}
                        <span className="font-light ml-1">{date?.getDate()}</span>
                    </p>
                    <div className="mt-4 flex flex-col h-[350px] bg-scroll gap-2 max-h-screen overflow-y-auto no-scrollbar text-center">
                        {(!isTimeAvailable() && date) && <> Currently not Available </>}
                        {(!isTimeAvailable() && !date) && <> Please select a date </>}
                        {(isTimeAvailable() && props?.eventType?.duration) &&
                            generateTimeRange(
                                startTime,
                                endTime,
                                props?.eventType?.duration
                            ).map((time, index) => (
                                // TODO: block date if in range was booked
                                // TODO: fix overtime
                                <Button
                                    variant="outline"
                                    key={index}
                                    onClick={() => setNextStep(time)}
                                >{time}</Button>
                            )
                        )}
                    </div>
                </div>
            </section>
        </AlertDialogContent>
    </AlertDialog>)
}