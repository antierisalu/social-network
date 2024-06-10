<script>
    import { onMount } from "svelte";

    let users = [];

    // Fetch users from the backend
    const fetchUsers = async () => {
        const response = await fetch('http://localhost:8080/allusers');
        if(response.ok) {
            const fetchedUsers = await response.json();
            users = [...fetchedUsers];
            console.log(users);
        } else {
            console.error('Error fetching users:', response.status);
        }
    };

    console.log(users)

    let searchQuery = '';

    // Reactive declaration for filtered users
    $: filteredUsers = searchQuery ? searchUsers() : users;

    const searchUsers = () => {
    return users.filter((user) => {
        return user.FirstName.toLowerCase().includes(searchQuery.toLowerCase()) ||
            user.LastName.toLowerCase().includes(searchQuery.toLowerCase())
    });
    };

    onMount(() => {
    fetchUsers();
  });

</script>


<div class="searchUsers">
    <input type="search" bind:value={searchQuery} placeholder="Search users...">
    {#if searchQuery}
        {#each filteredUsers as user}
            <div>
                <p>
                    <!-- svelte-ignore a11y-missing-attribute -->
                    <img src={user.Avatar}/> {user.FirstName} {user.LastName}
                </p>
            </div>
        {/each}
    {/if}
</div>



<style>

    input{
        margin: 0;
        border: none;
        width: 100%;
        height: 100%;
    }
  
    .searchUsers {
      margin-bottom: 8px;
    }

    img {
        max-height: 20px;
    }

</style>