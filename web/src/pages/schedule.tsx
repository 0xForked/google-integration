import {useParams} from "react-router-dom";

export function Schedule() {
    const {id} = useParams();
    return (<>{id}</>)
}