import {useEffect, useState} from "react";
import {
    AlertDialog,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogHeader,
    AlertDialogTitle
} from "@/components/ui/alert-dialog.tsx";
import {Globe} from "lucide-react";
import {Switch} from "@/components/ui/switch.tsx";
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectTrigger,
    SelectValue
} from "@/components/ui/select.tsx";
import {intToDay, intToTime, timeReference} from "@/lib/time.ts";
import {getAvailability} from "@/lib/api.ts";
import {Badge} from "@/components/ui/badge.tsx";
import {cn} from "@/lib/utils.ts";

interface AvailabilityDialogProps {
    display: boolean,
    callback(): void
}

export function AvailabilityModal(props: AvailabilityDialogProps) {
    const [display, setDisplay] = useState(false)
    const [availability, setAvailability] = useState<Availability | null>(null)
    const [availabilityFetched, setAvailabilityFetched] = useState(false);

    useEffect(() => {
        if (!availabilityFetched) {
            getUserAvailabilities();
        }
    }, []);

    const getUserAvailabilities = () => {
        getAvailability()
            .then((resp) => {
                setAvailability(resp.data);
                setAvailabilityFetched(true);
                setDisplay(props.display);
            })
            .catch((error) => {
                alert(error.message);
            });
    };

    useEffect(() => {
        if (availabilityFetched) {
            setDisplay(props.display);
        }
    }, [props.display, availabilityFetched]);

    return (<AlertDialog open={display} onOpenChange={props.callback}>
        <AlertDialogContent>
            <AlertDialogHeader>
                <AlertDialogTitle className="flex flex-row items-center gap-2">
                    {availability?.label}
                    <Badge>default</Badge>
                </AlertDialogTitle>
                <AlertDialogDescription>
                    Mon - Fri, 09:00 - 16:00
                    <span className="flex flex-row items-center mt-1">
                        <Globe className="w-4 h-4 inline mr-1" />
                        {availability?.timezone}
                    </span>
                </AlertDialogDescription>
            </AlertDialogHeader>

            <section className="flex flex-col my-4 gap-6">
                {availability?.days.map((ad, index) => (
                    <div className="grid grid-cols-2 items-center" key={index}>
                        <div className="flex flex-row gap-2 items-center ml-3">
                            <Switch disabled checked={!!ad.enable}/>
                            <span>{intToDay(ad.day)}</span>
                        </div>
                        <div className={cn(
                            !!ad.enable
                                ? "flex flex-row gap-2 items-center ml-3"
                                : "hidden"
                        )}>
                            <Select
                                defaultValue={!!ad.enable ? intToTime(ad.start_time) : ""}
                                onValueChange={() => alert("Sorry, updates are currently not supported, and this change will not be saved into the database.")}
                            >
                                <SelectTrigger className="w-[100px]">
                                    <SelectValue placeholder="Start" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectGroup>
                                        {timeReference.map((time, index) => (
                                            <SelectItem key={index} value={time}>{time}</SelectItem>
                                        ))}
                                    </SelectGroup>
                                </SelectContent>
                            </Select>
                            <span>-</span>
                            <Select
                                defaultValue={!!ad.enable ? intToTime(ad.end_time) : ""}
                                onValueChange={() => alert("Sorry, updates are currently not supported, and this change will not be saved into the database.")}
                            >
                                <SelectTrigger className="w-[100px]">
                                    <SelectValue placeholder="End" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectGroup >
                                        {timeReference.map((time, index) => (
                                            <SelectItem key={index} value={time}>{time}</SelectItem>
                                        ))}
                                    </SelectGroup>
                                </SelectContent>
                            </Select>
                        </div>
                    </div>
                ))}
            </section>

            <AlertDialogCancel>Close</AlertDialogCancel>
        </AlertDialogContent>
    </AlertDialog>)
}
