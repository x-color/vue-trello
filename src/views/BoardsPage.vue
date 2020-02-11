<template>
  <v-container>
    <!-- PIN: Title -->
    <v-row justify="start">
      <v-col cols="auto">
        <h1 class="ma-3">{{ user.name }}'s Boards</h1>
      </v-col>
    </v-row>

    <!-- PIN: Boards -->
    <draggable
      class="row flex-nowrap row--dense justify-start"
      group="boards"
      v-model="boards"
      draggable=".item"
    >
      <v-col v-for="(board, i) in boards" :key="i" cols="auto" class="item">
        <v-card dark>
          <router-link style="text-decoration: none" :to="`/board/${board.id}`">
            <Board :id="board.id" />
          </router-link>
        </v-card>
      </v-col>
      <v-col cols="auto">
        <v-btn
          class="mx-auto"
          width="300"
          height="100"
          color="grey lighten-2"
          @click="addBoardMode = true"
        >
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </v-col>
    </draggable>

    <!-- PIN: Modal for creating new board -->
    <BoardModal
      v-model="newBoard"
      :open="addBoardMode"
      @close="resetNewBoard"
      @save="addNewBoard"
    />
  </v-container>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import draggable from 'vuedraggable';
import BoardModal from '@/components/BoardModal.vue';
import Board from '@/components/Board.vue';

export default {
  name: 'boards-page',
  components: {
    Board,
    BoardModal,
    draggable,
  },
  computed: {
    ...mapGetters(['getBoardsByUserId']),
    user() {
      return this.$store.state.user.user;
    },
    boards: {
      get() {
        return this.getBoardsByUserId;
      },
      set(newValue) {
        const newUser = Object.assign(
          { ...this.user },
          { boards: newValue.map(board => board.id) },
        );
        this.editUser(newUser);
      },
    },
  },
  methods: {
    ...mapActions(['addBoard', 'editUser']),
    addNewBoard() {
      this.addBoard({
        userId: this.user.id,
        color: this.newBoard.color,
        text: this.newBoard.text,
        title: this.newBoard.title,
      });
      this.resetNewBoard();
    },
    resetNewBoard() {
      this.newBoard = { title: '', text: '', color: 'indigo' };
      this.addBoardMode = false;
    },
  },
  data() {
    return {
      addBoardMode: false,
      newBoard: { title: '', text: '', color: 'indigo' },
    };
  },
};
</script>
