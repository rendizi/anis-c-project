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
            //Если ИИН больше 12 или введенный элемент не цифра
            if (login.Text.Length > 12 || !char.IsDigit(login.Text[login.Text.Length - 1]))
            {
                //Вырезать новый добавленные невалидный элемент
                login.Text = login.Text.Substring(0, login.Text.Length - 1);
                login.SelectionStart = login.Text.Length;
            }
        }

        private void passwordShow_Click(object sender, EventArgs e)
        {
            //меняем иконку
            show = !show;
            string imagePath = "C:\\Users\\rendi\\RiderProjects\\aNIS\\aNIS\\imgs\\";

            //если пароль нужно показать
            if (show)
            {
                //меняем иконку на скрыть 
                passwordShow.Image = Image.FromFile(imagePath + "hide.png");
                password.PasswordChar = '\0';
            }
            else
            {
                //меняем иконку на показать
                passwordShow.Image = Image.FromFile(imagePath + "show.png");
                password.PasswordChar = '•';
            }
            password.UseSystemPasswordChar = !show;
        }

        private async void submit_Click(object sender, EventArgs e)
        {
            //ИИн должен быть 12
            if (login.Text.Length != 12)
            {
                MessageBox.Show("Invalid login");
                return;
            }

            //пароль не должен быть пустым
            if (password.Text == "")
            {
                MessageBox.Show("Password can't be an empty string");
                return;
            }

            //если ничего не выбрано- ошибка
            if (comboBox1.SelectedItem == null)
            {
                MessageBox.Show("Choose your school");
                return;
            }

            //если не согласен- ошибка
            if (!checkBox1.Checked)
            {
                MessageBox.Show("Agree with our politics");
                return;
            }
            
            try
            {
                //отправляем запрос на локальный сервер
                using (var client = new HttpClient())
                {
                    string url = $"http://localhost:8080/login?login={login.Text}&password={password.Text}";
                    var response = await client.GetAsync(url);

                    //проверяем статус реквеста
                    if (response.IsSuccessStatusCode)
                    {
                        //переносим на форму
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
            //перекидываем на мой сайт 
            System.Diagnostics.Process.Start("microsoft-edge:http://baglanov.site");
        }

        private void comboBox1_SelectedIndexChanged(object sender, EventArgs e)
        {
            
        }
    }
}
