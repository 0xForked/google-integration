import {useParams} from "react-router-dom";
import {Avatar, AvatarFallback} from "@/components/ui/avatar.tsx";
import {Badge} from "@/components/ui/badge.tsx";
import {Clock12} from "lucide-react";

export function Schedule() {
    const { username } = useParams();
    return (<div className="flex flex-col items-center py-12">
        <section className="mb-12">
            <Avatar className="mb-4 h-20 w-20">
                <AvatarFallback>
                    {username?.substring(1, 3).toUpperCase() ?? "-"}
                </AvatarFallback>
            </Avatar>
            <h1 className="text-lg font-bold">{username}</h1>
        </section>
        <section className="flex flex-col gap-2">
            <button
                className="w-96 h-20 border border-input bg-background hover:bg-accent hover:text-accent-foreground rounded-md p-4">
                <div className="flex flex-col text-left">
                    <h5 className="text-sm font-bold text-gray-600 mb-2">15 Min Meeting</h5>
                    <Badge variant="secondary" className="font-light w-16">
                        <Clock12 className="w-[12px] h-[12px] mr-1"/>15m
                    </Badge>
                </div>
            </button>
            <button
                className="w-96 h-20 border border-input bg-background hover:bg-accent hover:text-accent-foreground rounded-md p-4">
                <div className="flex flex-col text-left">
                    <h5 className="text-sm font-bold text-gray-600 mb-2">30 Min Meeting</h5>
                    <Badge variant="secondary" className="font-light w-16">
                        <Clock12 className="w-[12px] h-[12px] mr-1"/>30m
                    </Badge>
                </div>
            </button>
        </section>
    </div>)
}