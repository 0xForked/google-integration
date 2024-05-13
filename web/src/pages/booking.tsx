import {useNavigate, useParams} from "react-router-dom";
import {Avatar, AvatarFallback} from "@/components/ui/avatar.tsx";
import {Badge} from "@/components/ui/badge.tsx";
import {Clock12} from "lucide-react";
import {useEffect, useState} from "react";
import {getBookingHost} from "@/lib/api.ts";
import {BookingCalendarModal} from "@/components/booking-calendar-modal.tsx";
import {BookingFormModal} from "@/components/booking-form-modal.tsx";

export function Booking() {
    const { username } = useParams();
    const navigate = useNavigate()
    const [host, setHost] = useState<User | null>(null)
    const [selectedEventType, setSelectedEventType] = useState<EventType | null>(null)
    const [selectedDate, setSelectedDate] = useState<Date | undefined>(undefined)
    const [selectedTime, setSelectedTime] = useState<string | undefined>(undefined)
    const [displayCalendarDialog, setDisplayCalendarDialog] = useState(false)
    const [displayFormDialog, setDisplayFormDialog] = useState(false)

    useEffect(() => {
        if (!username) {
            if (confirm("username is required")) {
                window.location.reload()
            }
            return
        }
        getHost(username)
    }, [username])

    const getHost = (username: string) => {
        getBookingHost(username).then((resp) => {
            if (resp.error) {
                confirm(resp.error)
                navigate("/404")
                return
            }
            setHost(resp)
        }).catch((error) => alert(error.message))
    }

    const bookingForm = (date: Date, time: string)  => {
        setSelectedDate(date)
        setSelectedTime(time)
        setDisplayFormDialog(true)
    }

    return (<>
        <div className="flex flex-col items-center py-12">
            <section className="mb-12 flex flex-col text-center">
                <Avatar className="mb-4 h-20 w-20">
                    <AvatarFallback>
                        {host?.username?.substring(0, 2).toUpperCase() ?? "-"}
                    </AvatarFallback>
                </Avatar>
                <h1 className="text-lg font-bold">{host?.username.toUpperCase()}</h1>
            </section>
            <section className="flex flex-col gap-2">
                {host?.event_types?.map((et, index) => (
                    <button
                        key={index}
                        onClick={() => {
                            setSelectedEventType(et)
                            setDisplayCalendarDialog(true)
                        }}
                        className="w-96 h-20 border border-input bg-background hover:bg-accent hover:text-accent-foreground rounded-md p-4"
                    >
                        <div className="flex flex-col text-left">
                            <h5 className="text-sm font-bold text-gray-600 mb-2">{et.title}</h5>
                            <Badge variant="secondary" className="font-light w-16">
                                <Clock12 className="w-[12px] h-[12px] mr-1"/>{et.duration}m
                            </Badge>
                        </div>
                    </button>
                ))}
            </section>
        </div>

        <BookingCalendarModal
            username={host?.username}
            eventType={selectedEventType}
            display={displayCalendarDialog}
            callback={() => setDisplayCalendarDialog(false)}
            nextStep={bookingForm}
        />

        <BookingFormModal
            username={host?.username}
            eventType={selectedEventType}
            date={selectedDate}
            time={selectedTime}
            display={displayFormDialog}
            callback={() => setDisplayFormDialog(false)}
        />
    </>)
}