<template>
  <div class="register">
    <el-card>
      <h2>注册</h2>
      <el-form
          class="register-form"
      >
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="用户名" prefix-icon="fas fa-user"></el-input>
        </el-form-item>
        <el-form-item prop="email">
          <el-input v-model="form.email" placeholder="邮箱" prefix-icon="fas fa-user"></el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input
              v-model="form.password"
              placeholder="密码"
              type="password"
              prefix-icon="fas fa-lock"
          ></el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input
              v-model="form.repassword"
              placeholder="确认密码"
              type="password"
              prefix-icon="fas fa-lock"
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-button
              :loading="loading"
              class="login-button"
              type="primary"
              native-type="submit"
              @click="submitForm"
              block
          >注册</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import request from '../../utils/request'


export default {
  name: "register",
  data() {
    return {
      form: {
        username: '',
        password: '',
        repassword: '',
        email: '',
      },
      loading: false,
    };
  },
  methods: {
    submitForm: function(e) {
      if (this.form.repassword !== this.form.password) {
        this.$message.error("两次输入密码不一致");
        return;
      }
      const that = this;
      request({
        url: '/user/register',
        method: 'post',
        data: this.form
      }).then(function (response) {
        debugger
        const msg = response.data
        that.$message.success(msg);
      })
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.register {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 20px;
}

.login-button {
  width: 100%;
  margin-top: 40px;
}
.register-form {
  width: 390px;
}
.forgot-password {
  margin-top: 10px;
}
.send_verify_code{

}
</style>
