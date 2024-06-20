//when the value of a store changes, all components
// that are subscribed to that store will be changed
import { writable } from 'svelte/store'

//user is not logged in 
export const loggedIn = writable(false);

export const activeTab = writable('Profile')

// Profile searchbar
export const userProfileData = writable({})

// Profile editing
export const isEditingProfile = writable(false);

export const newAboutMeStore = writable('')
export const showOverlay = writable(false)

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

export const uploadImageStore = writable(null)

// Connected with WS Online users (user ID's)
export const onlineUserStore = writable([])