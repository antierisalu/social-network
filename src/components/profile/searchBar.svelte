<script>
    import { onMount } from "svelte";
    import { allUsers } from "../../stores"

    $: users = $allUsers;
    console.log("userid", users)
    
    // Fetch users from the backend

    var searchQuery = '';

    // Reactive declaration for filtered users
    $: filteredUsers = searchQuery ? searchUsers(searchQuery) : users;

    const searchUsers = (searchQuery) => {
    const [firstNameQuery, lastNameQuery] = searchQuery.toLowerCase().trim().split(' ');
    return users.filter((user) => {
      if (!firstNameQuery && !lastNameQuery){
        return
      }

      let firstNameMatch;
      let firstNameMatch2;
      let lastNameMatch;
      let lastNameMatch2;

      if (!lastNameQuery || (firstNameQuery && lastNameQuery)){
       firstNameMatch = user.FirstName.toLowerCase().includes(firstNameQuery);
       firstNameMatch2 = user.LastName.toLowerCase().includes(firstNameQuery);
      }
      if (!firstNameQuery || (firstNameQuery && lastNameQuery)){
       lastNameMatch =  user.LastName.toLowerCase().includes(lastNameQuery);
       lastNameMatch2 =  user.FirstName.toLowerCase().includes(lastNameQuery);
    }

        return firstNameMatch || lastNameMatch || firstNameMatch2 || lastNameMatch2;
    });
    }

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