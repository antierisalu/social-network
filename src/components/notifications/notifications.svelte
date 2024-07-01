<script>
    import { notifications } from '../../websocket.js';
    import { onMount } from 'svelte';

    let notificationList = [];

    onMount(() => {
        const unsubscribe = notifications.subscribe((newNotifications) => {
            notificationList = newNotifications;
        });

        return () => {
            unsubscribe();
        };
    });


    //vaga scuffed peaks vist backendis mingit teemat tegema et deletida vmdea
     function clearNotifications() {
        notificationList = [];
    }
</script>

<main>
    <h1>Notifications({notificationList.length})</h1>
    {#if notificationList.length > 0}
        <ul>
            {#each notificationList as notification}
                <li>{notification.data}</li>
            {/each}
        </ul>
        <button on:click={clearNotifications}>Clear Notifications</button>
    {/if}
    {#if notificationList.length === 0}
        <p>No notifications</p>
    {/if}
</main>

<style>
    main {
        padding: 20px;
        font-family: Arial, sans-serif;
    }

    h1 {
        font-size: 24px;
        margin-bottom: 10px;
    }

    ul {
        list-style-type: none;
        padding: 0;
    }

    li {
        background-color: #f4f4f4;
        padding: 10px;
        margin-bottom: 5px;
        border-radius: 4px;
    }

    p {
        color: #888;
    }
</style>
