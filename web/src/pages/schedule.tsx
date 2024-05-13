import {useParams} from "react-router-dom";
import {useEffect, useState} from "react";
import {getSchedule} from "@/lib/api.ts";

export function Schedule() {
    const {id} = useParams();
    const [schedule, setSchedule] = useState({});

    useEffect(() => {
        if (!id) {
            if (confirm("username is required")) {
                window.location.reload()
            }
            return
        }
        getSelectedSchedule(id)
    }, [id])

    const getSelectedSchedule = (id: string) => {
        getSchedule(id).then((resp) => {
            if (resp.error) {
                confirm(resp.error)
                return
            }
            setSchedule(resp)
        }).catch((error) => alert(error.message))
    }

    return (<>
        <pre>{JSON.stringify(schedule, null, 2)}</pre>
    </>)
}