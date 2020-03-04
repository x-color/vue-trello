<template>
  <v-container>
    <!-- PIN: Title -->
    <v-row justify="start">
      <v-col cols="auto">
        <h1 class="ma-3">{{ $store.state.user.user.name }}'s Boards</h1>
      </v-col>
    </v-row>
    <the-boards />
  </v-container>
</template>

<script>
import TheBoards from '@/components/TheBoards.vue';

export default {
  name: 'PageBoards',
  components: {
    TheBoards,
  },
};
</script>
