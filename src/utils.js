import { allUsers } from "./stores";

//backend genereerib uuid ja front end paneb clienti session cookie paika.
export function updateSessionToken(token, expire) {
    var dateTime = new Date();
    dateTime.setTime(dateTime.getTime() + expire * 60 * 60 * 1000);
    var expires = "expires=" + dateTime.toUTCString();
    document.cookie = "sessionToken=" + token.String + ";" + expires;
  }

export const fetchUsers = async () => {
    const response = await fetch('http://localhost:8080/allusers');
    if(response.ok) {
        const fetchedUsers = await response.json();
        allUsers.set([...fetchedUsers])
        console.log(allUsers)
    } else {
        console.error('Error fetching users:', response.status);
    }
};