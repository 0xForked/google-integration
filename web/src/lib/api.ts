export const API_URL = `http://localhost:8000/api/v1`

export const API_ENDPOINT = {
    USER: {
        LOGIN: `${API_URL}/login`,
        PROFILE: `${API_URL}/profile`,
        CALENDAR: `${API_URL}/profile/calendar`,
        LOGOUT: `${API_URL}/logout`,
    },
}

const signIn = async (username: string, password: string)  => {
    try {
        const response = await fetch(API_ENDPOINT.USER.LOGIN, {
            method: "POST",
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body: JSON.stringify({username, password})
        })
        const content = await response.json();
        return Promise.resolve(content)
    } catch (e) {
        return Promise.reject(e)
    }
}

const profile = async ()  => {
    try {
        const response = await fetch(API_ENDPOINT.USER.PROFILE, {
            method: "GET",
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
        })
        const content = await response.json();
        return Promise.resolve(content)
    } catch (e) {
        return Promise.reject(e)
    }
}

const calendar = async ()  => {
    try {
        const response = await fetch(API_ENDPOINT.USER.CALENDAR, {
            method: "GET",
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
        })
        const content = await response.json();
        return Promise.resolve(content)
    } catch (e) {
        return Promise.reject(e)
    }
}


const signOut = async () => {
    try {
        const response = await fetch(API_ENDPOINT.USER.LOGOUT, {
            method: "POST",
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
        })
        return Promise.resolve(response)
    } catch (e) {
        return Promise.reject(e)
    }
}

export  {
    signIn,
    signOut,
    profile,
    calendar
}