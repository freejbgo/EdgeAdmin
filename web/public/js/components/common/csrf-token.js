Vue.component("csrf-token", {
	created: function () {
		this.refreshToken()
	},
	mounted: function () {
		let that = this
		this.$refs.token.form.addEventListener("submit", function () {
			that.refreshToken()
		})
	},
	data: function () {
		return {
			token: ""
		}
	},
	methods: {
		refreshToken: function () {
			let that = this
			Tea.action("/csrf/token")
				.get()
				.success(function (resp) {
					that.token = resp.data.token
				})
		}
	},
	template: `<input type="hidden" name="csrfToken" :value="token" ref="token"/>`
})
