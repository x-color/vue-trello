<template>
  <v-container>
    <!-- PIN: Boards -->
    <draggable
      class="row flex-nowrap row--dense justify-start"
      group="boards"
      v-model="boards"
      draggable=".item"
      :animation="300"
    >
      <v-col v-for="(board, i) in boards" :key="i" cols="auto" class="item">
        <v-card dark>
          <router-link :to="`/boards/${board.id}`" class="router-no-underline">
            <card-board :id="board.id" />
          </router-link>
        </v-card>
      </v-col>
      <v-col cols="auto">
        <v-btn
          class="mx-auto"
          width="300"
          height="100"
          color="grey lighten-2"
          @click="openAddBoardModal"
        >
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </v-col>
    </draggable>

    <!-- PIN: Modal for creating new board -->
    <modal-board
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
import ModalBoard from '@/components/ModalBoard.vue';
import CardBoard from '@/components/CardBoard.vue';

export default {
  name: 'TheBoards',
  components: {
    CardBoard,
    ModalBoard,
    draggable,
  },
  computed: {
    ...mapGetters(['getSortedBoards']),
    boards: {
      get() {
        return this.getSortedBoards;
      },
      set(newValue) {
        this.moveBoard(newValue.map(board => board.id));
      },
    },
  },
  data() {
    return {
      addBoardMode: false,
      newBoard: { title: '', text: '', color: '' },
    };
  },
  created() {
    this.loadBoards(this.$store.state.user.user);
  },
  methods: {
    ...mapActions(['addBoard', 'loadBoards', 'moveBoard']),
    addNewBoard() {
      this.addBoard({
        userId: this.$store.state.user.user.id,
        color: this.newBoard.color,
        text: this.newBoard.text,
        title: this.newBoard.title,
      });
      this.resetNewBoard();
    },
    resetNewBoard() {
      this.newBoard = { title: '', text: '', color: this.$store.state.resource.colors[0] };
      this.addBoardMode = false;
    },
    openAddBoardModal() {
      this.resetNewBoard();
      this.addBoardMode = true;
    },
  },
};
</script>

<style scoped>
.router-no-underline {
  text-decoration: none;
}
</style>
