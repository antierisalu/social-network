<script>

    import Button from "../shared/button.svelte";

    let userData = [
        {
            email: "",
            firstName: "",
            lastName: "",
            dateOfBirth: "",
            avatar: "",
            nickName: "",
            aboutMe: "",
            password: "",
            passwordConfirm: "", 
        }
    ]

</script>

<div class="register">
    <form on:submit|preventDefault>
        <input type="email" placeholder="E-mail *" bind:value={userData.email} required>
        <input type="text" placeholder="First Name *" bind:value={userData.firstName} required>
        <input type="text" placeholder="Last Name *" bind:value={userData.lastName} required>
        <input type="date" bind:value={userData.dateOfBirth} required>
        <div class="fakeInput">
            <label for="avatar">Upload avatar (Optional)</label>
            <input type="file" id="avatar" class="hidden" accept="image/png, image/jpeg" bind:value={userData.avatar}>
        </div>
        <input type="text" placeholder="Nickname (Optional)" bind:value={userData.nickName}>
        <input type="text" placeholder="About Me (Optional)" bind:value={userData.aboutMe}>
        <input type="password" placeholder="Password" required bind:value={userData.password}>
        <input type="password" placeholder="Confirm Password" required bind:value={userData.passwordConfirm}>
        <Button type="secondary" on:click={console.table(userData)}>Register</Button>
    </form>
    <Button type="secondary" inverse={true} btn200px={true} on:click>Login Instead</Button>
</div>

<style>

    .fakeInput {
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
    }

    .register {
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    form {
        display: flex;
        flex-direction: column;
    }

    input{
        width: 300px;
        border-radius: 6px;
        padding: 8px 12px;
    }

    ::-webkit-calendar-picker-indicator {
    filter: invert(1);
    }
</style>