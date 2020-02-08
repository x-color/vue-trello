<template>
  <v-container>
    <v-row justify="start">
      <v-col cols="auto">
        <h1 class="ma-3">Your Boards</h1>
      </v-col>
    </v-row>

    <v-row dense justify="start">
      <v-col v-for="(board, i) in boards" :key="i" cols="auto">
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
    </v-row>

    <BoardModal
      v-model="newBoard"
      :open="addBoardMode"
      @close="resetNewBoard"
      @save="addNewBoard"
    />
  </v-container>
</template>

<script>
import { mapActions } from 'vuex';
import BoardModal from '@/components/BoardModal.vue';
import Board from '@/components/Board.vue';

export default {
  name: 'boards',
  components: {
    Board,
    BoardModal,
  },
  computed: {
    boards() {
      return this.$store.state.board.boards;
    },
  },
  methods: {
    ...mapActions(['addBoard']),
    addNewBoard() {
      this.addBoard({
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
