import { writable } from 'svelte/store'

export const loggedIn = writable(false);

export const activeTab = writable('Profile')

// Profile searchbar
export const userProfileData = writable({})

export const allUsers = writable([])

//this is client's own info
export const userInfo = writable({})

// Auth errors (login/register)
export const authError = writable('');
export function displayUserAuthError(errorStr) {
    authError.set(errorStr)
    setTimeout(() => {
      authError.set('')
    }, 3000);
}

// User list (right-side bar)
// export const userList = writable([])