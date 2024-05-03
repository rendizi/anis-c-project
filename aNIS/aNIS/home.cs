using System;
using System.Drawing;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;


namespace aNIS
{
    public partial class home : Form
    {
        private string token;
        private Form1 login;
        private string gid;
        public home(string responseBod, Form1 page)
        {
            InitializeComponent();
            token = responseBod;
            login = page;
            fetch(token);
        }

        private async void fetch(string token)
        {
            try
            {
                using (var client = new HttpClient())
                {
                    string url = $"http://localhost:8080/add?token={token}";
                    var response = await client.GetAsync(url);

                    if (response.IsSuccessStatusCode)
                    {
                        string responseBody = await response.Content.ReadAsStringAsync();
                        Console.WriteLine(responseBody);
                        string[] substrings = responseBody.Split(';');
                        label4.Text = substrings[1];
                        gid = substrings[0];
                    }
                    else
                    {
                        MessageBox.Show("Неверный логин или пароль");
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }

        private void label3_Click(object sender, EventArgs e)
        {
        }

        private void button2_Click(object sender, EventArgs e)
        {
            var BaPage = new b2a(this);
            this.Hide();
            BaPage.Show();
        }

        private void pictureBox2_Click_1(object sender, EventArgs e)
        {
            this.Hide();
            login.Show();
        }
         private void button1_Click_1(object sender, EventArgs e)
        {
            var GpaPage = new GPA(this, token, gid);
            this.Hide();
            GpaPage.Show();
        }

        private void label4_Click(object sender, EventArgs e)
        {
        }

       

        private void panel3_Paint(object sender, PaintEventArgs e)
        {
            
        }

        private void button2_Click_1(object sender, EventArgs e)
        {
            var page = new b2a(this);
            this.Hide();
            page.Show();
        }
    }
}