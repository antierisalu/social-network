<script>
    import Button from "../../shared/button.svelte";
    import { updateSessionToken, fetchUsers } from "../../utils";
    import { loggedIn, authError, displayUserAuthError} from "../../stores"
    import { fade, slide } from 'svelte/transition';
    import ImagePreview from "../../shared/imagePreview.svelte"
    let errorString = '';
    $: errorString = $authError;
    const passwordStrength = { PwLength: 5 }

    let userData = {
        email: "",
        firstName: "",
        lastName: "",
        dateOfBirth: "",
        avatar: "",
        avatarName: "",
        nickName: "",
        aboutMe: "",
        password: "",
        passwordConfirm: "",
    };

    // Send register info to backend
    async function registerUser(registerInfo) {
        console.log("Sending registerUser datato backend:", registerInfo)
        try {
        const response = await fetch("/register", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(registerInfo),
        });

        if (!response.ok) {
            const errorMessage = await response.text();
            throw new Error(errorMessage);
        }

        const data = await response.json();
        console.log("REGISTER:", data);
        updateSessionToken(data.session, 24);
        loggedIn.set(true)
        fetchUsers()
        } catch (error) {
            console.error("Error registering user:", error.message);
            displayUserAuthError(error.message)
        }
    }

    function checkPWstrength(password, passwordStrength) {
        if (password.length < passwordStrength.PwLength) {
            return false
        }

        return true
    }

    function isValidEmail(email) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    }

    // Convert dataURL to imgBlob
    function dataURLToBlob(dataURL) {
        const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png'];
        const parts = dataURL.split(',');
        const mime = parts[0].match(/:(.*?);/)[1];
        if (!allowedTypes.includes(mime)) {
            throw new Error('Invalid image format');
        }
        const bstr = atob(parts[1]);
        const n = bstr.length;
        const u8arr = new Uint8Array(n);
        for (let i = 0; i < n; i++) {
            u8arr[i] = bstr.charCodeAt(i);
        }

        return new Blob([u8arr], { type: mime });
    }

    // Convert imgBlob to base64
    function blobToBase64(blob) {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = () => resolve(reader.result.split(',')[1]);
            reader.onerror = reject;
            reader.readAsDataURL(blob);
        });
    }

    async function attemptRegister(userData) {
        console.table(userData)

        // 1. Check for empty fields
        if (
            userData.email === "" ||
            userData.firstName === "" ||
            userData.lastName === "" ||
            userData.dateOfBirth === "" ||
            userData.password === "" ||
            userData.passwordConfirm === ""
        ) {
            displayUserAuthError("Please fill in all the fields")
            return
        }
        // 2. Email validity
        if (!isValidEmail(userData.email)) {
            displayUserAuthError("Please use a valid email address")
            return
        }
        // 3. Password strength
        if (!checkPWstrength(userData.password, passwordStrength)) {
            displayUserAuthError("Please use a stronger password")
            return
        }
        // 4. Password match
        if (userData.password !== userData.passwordConfirm) {
            displayUserAuthError("Passwords dont match")
            return
        }
        // 5. Check if Avatar img exists
        const avatarPreview = document.getElementById('avatarPreview')
        if (avatarPreview !== null) {
            try {
                const imageDataURL = avatarPreview.src;
                const imgBlob = dataURLToBlob(imageDataURL)
                // Convert Blob to Base64
                try {
                    const base64String = await blobToBase64(imgBlob);
                    userData.avatar = base64String;
                    userData.avatarName = avatarPreview.getAttribute('name');
                } catch (error) {
                    console.error("Error converting blob to base64:", error);
                }
            } catch (error) {
                console.error('dataURLtoBlob catch error: ', error.message)
                displayUserAuthError(error.message)
                return
            }
        }

        registerUser(userData)
    }
</script>

<div class="register" in:fade>
    <form on:submit|preventDefault>
        <input type="email" placeholder="E-mail *" bind:value={userData.email} required>
        <input type="text" placeholder="First Name *" bind:value={userData.firstName} required>
        <input type="text" placeholder="Last Name *" bind:value={userData.lastName} required>
        <input type="date" bind:value={userData.dateOfBirth} required>
        <input type="password" placeholder="Password *" required bind:value={userData.password}>
        <input type="password" placeholder="Confirm Password *" required bind:value={userData.passwordConfirm}>
        <input type="text" placeholder="Nickname (Optional)" bind:value={userData.nickName}>
        <input type="text" placeholder="About Me (Optional)" bind:value={userData.aboutMe}>
        <ImagePreview />
        <Button type="secondary" on:click={attemptRegister(userData)}>Register</Button>
    </form>
        {#if errorString != ""}
        <div class="error" transition:slide>
            <Button type="primary" customStyle="width:300px; min-height: 35px; cursor: default; pointer-events: none;">{errorString}</Button>
        </div>
        {/if}
    <Button type="secondary" inverse={true} customStyle="width:200px" on:click>Login Instead</Button>
</div>

<style>
    /* .fakeInput {
    color: #ddd;
    display: inline-block;
    padding: 8px 10px;
    border: 1px solid #ccc;
    border-radius: 6px;
    background-color: #011;
    margin-bottom: 8px;
    cursor: pointer;
    }

    .fakeInput label {
    cursor: pointer;
    text-align: left;
    margin: 1;
    opacity: 0.7;
    }

    .hidden {
    display: none;
    } */

    .register {
    display: flex;
    flex-direction: column;
    align-items: center;
    }

    form {
    display: flex;
    flex-direction: column;
    }

    input {
    width: 300px;
    border-radius: 6px;
    padding: 8px 12px;
    }

    ::-webkit-calendar-picker-indicator {
    filter: invert(1);
    }
</style>
