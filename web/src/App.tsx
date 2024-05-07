import {LoginDialog} from "@/components/login-modal.tsx";
import {useEffect, useState} from "react";
import {calendar, profile, signOut} from "@/lib/api.ts";
import {Button} from "@/components/ui/button.tsx";

function App() {
    const [
        displayLoginDialog,
        setDisplayLoginDialog,
    ] = useState(false)
    const [
        user,
        setUser,
    ] = useState<User | null>(null)
    const [
        authUrl,
        setAuthUrl,
    ] = useState<string | null>(null)
    const [
        oAuthProfileName,
        setOAuthProfileName,
    ] = useState<string | null>(null)

    useEffect(() => {
        getUser()
    }, [])

    const getUser = () => {
        profile().then((resp) => {
            if (typeof resp === "string" && resp === "ACCESS_TOKEN_NOT_PROVIDE") {
                setDisplayLoginDialog(true)
                return
            }
            setUser(resp.data)
            getCalendar()
        }).catch((error) => alert(error.message))
    }

    const getCalendar = () => {
        calendar().then((resp) => {
            if (resp.auth_url != "" ) {
                setAuthUrl(resp.auth_url)
            }
            if (resp.name != "" ) {
                setOAuthProfileName(resp.name)
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

    return (
        <>
            <div className="md:hidden text-center py-10">
                Screen Size Not Supported <br/> (min: 768px/tablet screen)
            </div>
            <div className="hidden flex-col md:flex">
                <main className="min-h-full min-w-full p-10">
                    <section className="border-2 p-4 mb-4">
                        <h1 className="text-xl font-bold">Internal User Data</h1>
                        <h5 className="text-lg">{user?.username} </h5>
                        <Button
                            className="block mt-4"
                            size="sm"
                            onClick={() => logout()}
                        >Logout</Button>
                    </section>

                    <section className="border-2 p-4 mb-4">
                        <h1 className="text-xl font-bold">External User Data (Google OAuth)</h1>
                        {authUrl && <Button
                            onClick={() => openInNewTab(authUrl)}
                        >Connect with Google</Button>}
                        <h5 className="text-lg">{oAuthProfileName && oAuthProfileName}</h5>
                    </section>

                    <section className="border-2 p-4">
                        <h1 className="text-xl font-bold">Calendar Data (Google)</h1>
                        test . . .
                    </section>
                </main>
            </div>
            <LoginDialog display={displayLoginDialog}/>
        </>
    )
}

export default App
