using System;
using System.Drawing;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace aNIS
{
    public partial class Form1 : Form
    {
        public Form1()
        {
            InitializeComponent();
            comboBox1.DropDownStyle = ComboBoxStyle.DropDownList;

        }

        private bool show = false;

        private void login_TextChanged(object sender, EventArgs e)
        {
            if (login.Text.Length > 12 || !char.IsDigit(login.Text[login.Text.Length - 1]))
            {
                login.Text = login.Text.Substring(0, login.Text.Length - 1);
                login.SelectionStart = login.Text.Length;
            }
        }

        private void passwordShow_Click(object sender, EventArgs e)
        {
            show = !show;
            string imagePath = "C:\\Users\\rendi\\RiderProjects\\aNIS\\aNIS\\imgs\\";

            if (show)
            {
                passwordShow.Image = Image.FromFile(imagePath + "hide.png");
                password.PasswordChar = '\0';
            }
            else
            {
                passwordShow.Image = Image.FromFile(imagePath + "show.png");
                password.PasswordChar = '•';
            }
            password.UseSystemPasswordChar = !show;
        }

        private async void submit_Click(object sender, EventArgs e)
        {
            if (login.Text.Length != 12)
            {
                MessageBox.Show("Invalid login");
                return;
            }

            if (password.Text == "")
            {
                MessageBox.Show("Password can't be an empty string");
                return;
            }

            if (comboBox1.SelectedItem == null)
            {
                MessageBox.Show("Choose your school");
                return;
            }

            if (!checkBox1.Checked)
            {
                MessageBox.Show("Agree with our politics");
                return;
            }
            
            try
            {
                using (var client = new HttpClient())
                {
                    string url = $"http://localhost:8080/login?login={login.Text}&password={password.Text}";
                    var response = await client.GetAsync(url);

                    if (response.IsSuccessStatusCode)
                    {
                        string responseBody = await response.Content.ReadAsStringAsync();
                        var homeForm = new home(responseBody, this);
                        this.Hide();
                        homeForm.Show();
                    }
                    else
                    {
                        MessageBox.Show("Login failed. Please check your credentials.");
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show("An error occurred: " + ex.Message);
            }
        }


        private void password_TextChanged(object sender, EventArgs e)
        {
            if (password.Text.Length >= 1)
            {
                passwordShow.Show();
            }
            else
            {
                passwordShow.Hide();
            }
        }

        private void linkLabel1_LinkClicked(object sender, LinkLabelLinkClickedEventArgs e)
        {
            System.Diagnostics.Process.Start("microsoft-edge:http://baglanov.site");
        }

        private void comboBox1_SelectedIndexChanged(object sender, EventArgs e)
        {
            
        }
    }
}
