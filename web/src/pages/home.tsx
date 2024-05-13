import {LoginDialog} from "@/components/login-modal.tsx";
import {useEffect, useState} from "react";
import {getEvent, getProfile, signOut} from "@/lib/api.ts";
import {Button} from "@/components/ui/button.tsx";
import {
    DropdownMenu,
    DropdownMenuContent, DropdownMenuItem,
    DropdownMenuLabel, DropdownMenuSeparator,
    DropdownMenuTrigger
} from "@/components/ui/dropdown-menu.tsx";
import {AvailabilityModal} from "@/components/availability-modal.tsx";
import {EventTypeModal} from "@/components/event-modal.tsx";
import {Accordion, AccordionContent, AccordionItem, AccordionTrigger} from "@/components/ui/accordion.tsx";

function Home() {
    const [displayLoginDialog, setDisplayLoginDialog] = useState(false)
    const [displayAvailabilityDialog, setDisplayAvailabilityDialog] = useState(false)
    const [displayEventTypeDialog, setDisplayEventTypeDialog] = useState(false)
    const [user, setUser] = useState<User | null>(null)
    const [authUrl, setAuthUrl] = useState<string | null>(null)
    const [oAuthDisplayName, setOAuthDisplayName] = useState<string | null>(null)
    const [oAuthEmail, setOAuthEmail] = useState<string | null>(null)
    const [scheduledEvents, setScheduledEvents] = useState([]);

    useEffect(() => getUser(), [])

    const getUser = () => {
        getProfile().then((resp) => {
            if (typeof resp === "string" && resp === "ACCESS_TOKEN_NOT_PROVIDE") {
                setDisplayLoginDialog(true)
                return
            }
            setUser(resp.data)
            getCalendar()
        }).catch((error) => alert(error.message))
    }

    const getCalendar = () => {
        getEvent().then((resp) => {
            if (resp.auth_url != "" ) {
                setAuthUrl(resp.auth_url)
            }
            if (resp.name != "" || resp.email != "") {
                setOAuthDisplayName(resp.name)
                setOAuthEmail(resp.email)
            }
            if (resp.scheduled) {
                setScheduledEvents(resp.scheduled)
            }
        }).catch((error) => alert(error))
    }

    const openInNewTab = (url?: string| null) => {
        if (url == null) {
            alert("auth url is required")
            return
        }
        window.open(url, "_self", "noreferrer");
    };

    const logout = () => {
        signOut().then(() => {
            window.location.reload();
        }).catch((error) => alert(error))
    }

    const closeAvailabilityModal = () => setDisplayAvailabilityDialog(false)

    const closeEventTypeModal = () => setDisplayEventTypeDialog(false)

    return (<>
        <div className="md:hidden text-center py-10">
            Screen Size Not Supported <br/> (min: 768px/tablet screen)
        </div>
        <div className="hidden flex-col md:flex">
            <main className="min-h-full min-w-full p-10">
                {!user && <>
                    <div className="text-center py-10">
                        Please Login to Continue
                    </div>
                </>}
                {user && <>
                    <section className="border-2 p-4 mb-4 grid grid-cols-2">
                        <div className="grid">
                            <h1 className="text-xl font-bold">Internal User Data</h1>
                            <h5 className="text-lg">
                                @{user?.username}
                            </h5>
                            {user && <div className="flex gap-2 mt-4">
                                <DropdownMenu>
                                    <DropdownMenuTrigger className="rounded-md text-sm font-medium h-10 px-4 py-2 border border-input bg-background hover:bg-accent hover:text-accent-foreground">
                                        Update
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent>
                                        <DropdownMenuLabel>My Account</DropdownMenuLabel>
                                        <DropdownMenuSeparator/>
                                        <DropdownMenuItem
                                            onClick={() => setDisplayAvailabilityDialog(true)}>Availability</DropdownMenuItem>
                                        <DropdownMenuItem onClick={() => setDisplayEventTypeDialog(true)}>Events
                                            Types</DropdownMenuItem>
                                    </DropdownMenuContent>
                                </DropdownMenu>
                                <Button
                                    size="sm"
                                    variant="destructive"
                                    onClick={() => logout()}
                                >Logout</Button>
                            </div>}
                        </div>
                        <div>
                            <h1 className="text-xl font-bold">External User Data (Google OAuth)</h1>
                            {authUrl && <Button
                                onClick={() => openInNewTab(authUrl)}
                            >Connect with Google</Button>}
                            <h5 className="text-lg">{oAuthDisplayName}</h5>
                            <h5 className="text-lg">{oAuthEmail}</h5>
                        </div>
                    </section>

                    <section className="border-2 p-4">
                        <h1 className="text-xl font-bold">Events</h1>
                        <hr className="my-4"/>
                        <div className="grid grid-cols-2 divide-x-[1px]">
                            <div className="px-4 flex flex-col gap-2">
                                <h5 className="font-bold">
                                    Scheduled
                                </h5>
                                <Accordion type="single" collapsible>
                                    {scheduledEvents.map((item: any, index) => (
                                        <AccordionItem className="border-b" value={item?.id} key={index}>
                                            <AccordionTrigger>
                                                <div className="flex flex-col items-center gap-2">
                                                    <h5 className="text-sm font-bold text-gray-600">
                                                        {item?.summary}
                                                    </h5>
                                                    <p className="ml-[-125px] text-xs font-light">
                                                        {item?.hangoutLink}
                                                    </p>
                                                </div>
                                            </AccordionTrigger>
                                            <AccordionContent className="flex flex-col">
                                                <pre>{JSON.stringify(item, null, 2)}</pre>
                                            </AccordionContent>
                                        </AccordionItem>
                                    ))}
                                </Accordion>
                            </div>
                            <div className="px-4 flex flex-col gap-2">
                                <h5 className="font-bold">
                                    Canceled
                                </h5>
                            </div>
                        </div>
                    </section>
                </>}
            </main>
        </div>

        <LoginDialog display={displayLoginDialog}/>
        <AvailabilityModal display={displayAvailabilityDialog} callback={closeAvailabilityModal}/>
        <EventTypeModal display={displayEventTypeDialog} callback={closeEventTypeModal}/>
    </>)
}

export default Home
