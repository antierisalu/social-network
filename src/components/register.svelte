<script>

    let userData = [
        {
            email: email,
            firstName: firstName,
            lastName: lastName,
            

        }
    ]

    
</script>



<div>
    <form on:submit|preventDefault>
        <input type="text" placeholder="email" bind:value={email} required>
        <input type="text" placeholder="firstName" bind:value={firstName} required>
        <input type="text" placeholder="lastName" bind:value={lastName} required>
        <input type="date" placeholder="Date of birth" bind:value={dateOfBirth} required>
        <input type="text" placeholder="avatar" bind:value={avatar}>
        <input type="text" placeholder="nickName" bind:value={nickName}>
        <input type="text" placeholder="About Me" bind:value={aboutMe}>
        <input type="password" placeholder="password" required bind:value={password}>
        <Button type="secondary" on:click={fetchData}>Login</Button>
        
    </form>
</div>

<style>

    input{
        border-radius: 6px;
        padding: 8px 12px;
    }
</style>