<template>
  <div>
    <div v-if="board">
      <!-- PIN: Board info -->
      <v-container>
        <v-row>
          <h1 class="text-center ma-3">{{ board.title }}</h1>
          <div class="text-center my-auto">
            <v-menu offset-y>
              <template v-slot:activator="{ on }">
                <v-btn outlined small text dark v-on="on">
                  <v-icon color="black">mdi-dots-horizontal</v-icon>
                </v-btn>
              </template>
              <v-list>
                <v-list-item
                  class="pr-1"
                  v-for="(item, index) in menuItems"
                  :key="index"
                  @click="item.action()"
                >
                  <v-list-item-title>{{ item.title }}</v-list-item-title>
                  <v-list-item-avatar class="mx-0">
                    <v-icon small>{{ item.icon }}</v-icon>
                  </v-list-item-avatar>
                </v-list-item>
              </v-list>
            </v-menu>
          </div>
        </v-row>
        <v-row class="board-text">
          <p class="mx-4">{{ board.text }}</p>
        </v-row>

        <TheLists :id="board.id"/>

        <!-- PIN: Modal for editing board -->
        <ModalBoard
          v-model="editedBoard"
          :open="editBoardMode"
          @close="editBoardMode = false"
          @save="saveEditedBoard"
        />
      </v-container>

      <!-- PIN: Confirmation modal for deleting board -->
      <ModalConfirm
        :title="`Delete '${board.title}' ?`"
        text="Can not restore all data in board."
        :open="deleteBoardMode"
        @cancel="deleteBoardMode = false"
        @confirm="del()"
      />
    </div>

    <!-- PIN: No board page -->
    <div v-else>
      <v-container>
        <v-row class="text-center" justify="center">
          <v-col cols="12">
            <h1>This Board does not exist.</h1>
          </v-col>
          <v-col>
            <p>
              Please go to
              <router-link to="/boards">Home</router-link>
            </p>
          </v-col>
        </v-row>
      </v-container>
    </div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import ModalBoard from '@/components/ModalBoard.vue';
import ModalConfirm from '@/components/ModalConfirm.vue';
import TheLists from '@/components/TheLists.vue';

export default {
  name: 'PageBoard',
  components: {
    ModalBoard,
    ModalConfirm,
    TheLists,
  },
  computed: {
    ...mapGetters(['getBoardById']),
    board() {
      return this.getBoardById(this.$route.params.id);
    },
  },
  data() {
    return {
      editBoardMode: false,
      deleteBoardMode: false,
      editedBoard: {},
      menuItems: [
        {
          title: 'edit',
          icon: 'mdi-pencil',
          action: () => {
            this.editedBoard = { ...this.board };
            this.editBoardMode = true;
          },
        },
        {
          title: 'delete',
          icon: 'mdi-trash-can-outline',
          action: () => {
            this.deleteBoardMode = true;
          },
        },
      ],
    };
  },
  watch: {
    'board.color': {
      handler(newValue) {
        this.changeColor({ color: newValue });
      },
      immediate: true,
    },
  },
  methods: {
    ...mapActions(['addList', 'editBoard', 'removeBoard', 'editList', 'changeColor']),
    saveEditedBoard() {
      this.editBoard(this.editedBoard);
      this.editBoardMode = false;
    },
    del() {
      this.deleteBoardMode = false;
      this.removeBoard(this.board);
    },
  },
};
</script>

<style scoped>
.board-text {
  height: 15px;
}
</style>
