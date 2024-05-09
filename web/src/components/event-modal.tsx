import {useEffect, useState} from "react";
import {
    AlertDialog, AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogHeader,
    AlertDialogTitle
} from "@/components/ui/alert-dialog.tsx";
import {
    Accordion,
    AccordionContent,
    AccordionItem,
    AccordionTrigger,
} from "@/components/ui/accordion"
import {Badge} from "@/components/ui/badge.tsx";
import {Clock12} from "lucide-react";
import {getEventType} from "@/lib/api.ts";
import {intToDay, intToTime} from "@/lib/time.ts";

interface EventTypeDialogProps {
    display: boolean,
    callback(): void
}

export function EventTypeModal(props: EventTypeDialogProps) {
    const [display, setDisplay] = useState(false)
    const [eventTypes, setEventTypes] = useState<EventType[] | null>(null)
    const [eventTypeFetched, setEventTypeFetched] = useState(false);

    useEffect(() => {
        if (!eventTypeFetched) {
            getUserEventTypes();
        }
    }, []);

    const getUserEventTypes = () => {
        getEventType()
            .then((resp) => {
                setEventTypes(resp.data);
                setEventTypeFetched(true);
                setDisplay(props.display);
            })
            .catch((error) => {
                alert(error.message);
            });
    };

    useEffect(() => {
        if (eventTypeFetched) {
            setDisplay(props.display);
        }
    }, [props.display, eventTypeFetched]);

    return (<AlertDialog open={display} onOpenChange={props.callback}>
        <AlertDialogContent>
            <AlertDialogHeader>
                <AlertDialogTitle>Event Types</AlertDialogTitle>
                <AlertDialogDescription>
                    Create events to share for people to book on your calendar.
                    <span className="block text-xs text-yellow-500 mt-2">Updates are currently not supported</span>
                </AlertDialogDescription>
            </AlertDialogHeader>

            <section className="flex flex-col my-4">
                <Accordion type="single" collapsible>
                    {eventTypes?.map((et, index) => (
                        <AccordionItem className="border-b" value={`${et.duration}m`} key={index}>
                            <AccordionTrigger>
                                <div className="flex flex-row items-center gap-2">
                                    <h5 className="text-sm font-bold text-gray-600">
                                        {et.title}
                                    </h5>
                                    <Badge variant="secondary" className="font-light w-16">
                                        <Clock12 className="w-[12px] h-[12px] mr-1"/>{`${et.duration}m`}
                                    </Badge>
                                </div>
                            </AccordionTrigger>
                            <AccordionContent className="flex flex-col">
                                <h5 className="text-md font-bold">Title:</h5>
                                <span>
                                    {et.title}
                                    {!!et.enable && <Badge className="ml-1 h-[22px]">
                                        Enabled
                                    </Badge>}
                                </span>

                                <h5 className="text-md font-bold mt-2">Description:</h5>
                                <p>{et.description}</p>

                                <h5 className="text-md font-bold mt-2">Availability:</h5>
                                <span className="font-semibold">
                                    {et.availability.label}
                                    {!!et.enable && <Badge variant="outline" className="ml-1 h-[22px]">
                                        Default
                                    </Badge>}
                                </span>
                                {et.availability.days.map((ed, index) => (
                                    <div key={index} className="grid grid-cols-2">
                                        <p className={!!ed.enable ? "" : "line-through"}>{intToDay(ed.day)}</p>
                                        <span className="flex flex-row text-gray-500 items-center">
                                            {!!ed.enable
                                                ? `${intToTime(ed.start_time)} - ${intToTime(ed.end_time)}`
                                                : "unavailable"
                                            }
                                        </span>
                                    </div>
                                ))}

                                <h5 className="text-md font-bold mt-2">Timezone:</h5>
                                <p>{et.availability.timezone}</p>

                                <h5 className="text-md font-bold mt-2">Duration:</h5>
                                <p>{et.duration} Minutes</p>

                                <h5 className="text-md font-bold mt-2">Location:</h5>
                                <p>Google Meet</p>
                            </AccordionContent>
                        </AccordionItem>
                    ))}
                </Accordion>
            </section>

            <AlertDialogCancel>Close</AlertDialogCancel>
        </AlertDialogContent>
    </AlertDialog>)
}