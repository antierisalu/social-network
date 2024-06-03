<script>
    import Button from "../shared/button.svelte";
    import { updateSessionToken } from "../utils";
    import { loggedIn } from "../stores"
    import { fade, slide } from 'svelte/transition';
    import ImagePreview from "../shared/imagePreview.svelte"
    $: errorString = "";
    const passwordStrength = { PwLength: 5 }

    let userData = {
        email: "",
        firstName: "",
        lastName: "",
        dateOfBirth: "",
        avatar: "",
        nickName: "",
        aboutMe: "",
        password: "",
        passwordConfirm: "",
    };

    //send register info to backend
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
        updateSessionToken(data.token, data.expires);
        loggedIn.set(true)
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

    // Display error
    function displayUserAuthError(errorStr) {
        errorString = `${errorStr}`
        setTimeout(() => {
            errorString = ``
        }, 3000);
    }

    function isValidEmail(email) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    }

    // function getImage

    function attemptRegister(userData) {
        console.table(userData)
        // 5. Check if Avatar img exists
        const avatarPreview = document.getElementById('avatarPreview')
        if (avatarPreview !== null) {
            // console.log("avatar preview:", avatarPreview.src)

            // Conv to Img blob
            const imageDataURL = avatarPreview.src;
            const parts = imageDataURL.split(',');
            const mime = parts[0].match(/:(.*?);/)[1];
            const bstr = atob(parts[1]);
            let n = bstr.length;
            let u8arr = new Uint8Array(n);
            while(n--) {
                u8arr[n] = bstr.charCodeAt(n);
            }

            const imgBlob = new Blob([u8arr], {type: mime});
            console.log(imgBlob);
            
            userData.avatar = 'imgBLOBHERE'
        }

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

        // // 5. Check if Avatar img exists (MOVED UP FOR TESTING)
        // const avatarPreview = document.getElementById('avatarPreview')
        // if (avatarPreview !== null) {
        //     console.log("avatar preview:", avatarPreview)
        // }

        
        // Attempt user-creation (backend<errors>)
        registerUser(userData)
        // 1. Email taken
        // 2. Nickname taken
        // 3. Name+LastName taken
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
