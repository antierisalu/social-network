<script>
    import Button from "../shared/button.svelte";
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();
    export function logout() {
        dispatch('logout',{
            loginStatus: false
        })
    }
</script>

    <Button type="secondary" btn200px={true} inverse={true}  on:click={logout}>Logout</Button>


<style>


</style>